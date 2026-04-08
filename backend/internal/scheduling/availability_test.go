package scheduling

import (
	"testing"
	"time"
)

// ─── Helpers ──────────────────────────────────────────────────────────────────

const testTimezone = "America/Sao_Paulo"

// makeDate cria um time.Time representando "YYYY-MM-DD 00:00:00 -03:00 (BRT)"
// convertido para UTC.
func makeDate(year int, month time.Month, day int) time.Time {
	loc, _ := time.LoadLocation(testTimezone)
	return time.Date(year, month, day, 0, 0, 0, 0, loc).UTC()
}

// makeDateTime cria um time.Time em UTC a partir de hora/minuto no fuso BRT.
func makeDateTime(year int, month time.Month, day, hour, minute int) time.Time {
	loc, _ := time.LoadLocation(testTimezone)
	return time.Date(year, month, day, hour, minute, 0, 0, loc).UTC()
}

// bizHours cria um slice com um único dia de funcionamento aberto.
func bizHoursOpen(dow int, open, close string) []BusinessHourRow {
	return []BusinessHourRow{
		{DayOfWeek: dow, OpenTime: open, CloseTime: close, IsClosed: false},
	}
}

// profHoursAvailable cria um slice com um único dia disponível para o profissional.
func profHoursAvailable(dow int, start, end string) []ProfessionalHourRow {
	return []ProfessionalHourRow{
		{DayOfWeek: dow, StartTime: start, EndTime: end, IsUnavailable: false},
	}
}

// ─── Cenário 1: slot disponível normal ────────────────────────────────────────

func TestCalculateSlots_SlotLivre(t *testing.T) {
	// Segunda-feira 2026-04-06 (dow=1)
	// Estabelecimento: 09:00–18:00
	// Profissional: 09:00–18:00 disponível
	// Serviço: 30 minutos
	// Espera: vários slots, o primeiro às 09:00 e o último às 17:30

	date := makeDate(2026, time.April, 6)
	from := date
	to := date.Add(24 * time.Hour)

	dow := int(date.In(mustLoadLocation(testTimezone)).Weekday())

	biz := bizHoursOpen(dow, "09:00:00", "18:00:00")
	prof := profHoursAvailable(dow, "09:00:00", "18:00:00")

	slots := CalculateSlots(
		biz, prof,
		nil, nil, nil,
		30, from, to,
		testTimezone, "prof-1",
	)

	if len(slots) == 0 {
		t.Fatal("esperava slots disponíveis, mas nenhum foi retornado")
	}

	firstSlot := slots[0]
	loc := mustLoadLocation(testTimezone)

	expectedStart := makeDateTime(2026, time.April, 6, 9, 0)
	if !firstSlot.StartsAt.Equal(expectedStart) {
		t.Errorf("primeiro slot: esperava starts_at=%v, got=%v", expectedStart.In(loc), firstSlot.StartsAt.In(loc))
	}

	lastSlot := slots[len(slots)-1]
	expectedLastStart := makeDateTime(2026, time.April, 6, 17, 30)
	if !lastSlot.StartsAt.Equal(expectedLastStart) {
		t.Errorf("último slot: esperava starts_at=%v, got=%v", expectedLastStart.In(loc), lastSlot.StartsAt.In(loc))
	}

	// 9h às 18h = 9h = 540 min / 30 min = 18 slots
	if len(slots) != 18 {
		t.Errorf("esperava 18 slots, got %d", len(slots))
	}
}

// ─── Cenário 2: slot bloqueado por agendamento ───────────────────────────────

func TestCalculateSlots_SlotOcupadoAppointment(t *testing.T) {
	date := makeDate(2026, time.April, 6)
	from := date
	to := date.Add(24 * time.Hour)

	dow := int(date.In(mustLoadLocation(testTimezone)).Weekday())

	biz := bizHoursOpen(dow, "09:00:00", "18:00:00")
	prof := profHoursAvailable(dow, "09:00:00", "18:00:00")

	// Agendamento que ocupa 09:00–09:30
	apptStart := makeDateTime(2026, time.April, 6, 9, 0)
	apptEnd := makeDateTime(2026, time.April, 6, 9, 30)
	appointments := []Appointment{
		{
			ID:             "appt-1",
			ProfessionalID: "prof-1",
			StartsAt:       apptStart,
			EndsAt:         apptEnd,
			Status:         "confirmed",
		},
	}

	slots := CalculateSlots(
		biz, prof,
		appointments, nil, nil,
		30, from, to,
		testTimezone, "prof-1",
	)

	// O slot das 09:00 não deve aparecer.
	for _, s := range slots {
		if s.StartsAt.Equal(apptStart) {
			t.Error("slot das 09:00 não deveria estar disponível (ocupado por agendamento)")
		}
	}

	// Deve ter 17 slots (18 - 1 bloqueado).
	if len(slots) != 17 {
		t.Errorf("esperava 17 slots, got %d", len(slots))
	}
}

// ─── Cenário 3: fora do horário de funcionamento ─────────────────────────────

func TestCalculateSlots_ForaDoHorario(t *testing.T) {
	// Estabelecimento abre às 14:00; serviço de 30 min.
	// Solicitamos slots a partir de 09:00 — não deve haver nada antes das 14:00.
	date := makeDate(2026, time.April, 6)
	from := date
	to := date.Add(24 * time.Hour)

	dow := int(date.In(mustLoadLocation(testTimezone)).Weekday())

	biz := bizHoursOpen(dow, "14:00:00", "18:00:00")
	prof := profHoursAvailable(dow, "09:00:00", "18:00:00")

	slots := CalculateSlots(
		biz, prof,
		nil, nil, nil,
		30, from, to,
		testTimezone, "prof-1",
	)

	loc := mustLoadLocation(testTimezone)
	for _, s := range slots {
		localStart := s.StartsAt.In(loc)
		if localStart.Hour() < 14 {
			t.Errorf("slot às %v está antes do horário de abertura (14:00)", localStart)
		}
	}

	// 14h às 18h = 4h = 240 min / 30 min = 8 slots
	if len(slots) != 8 {
		t.Errorf("esperava 8 slots, got %d", len(slots))
	}
}

// ─── Cenário 4: duração não cabe no último slot ──────────────────────────────

func TestCalculateSlots_DuracaoNaoCabe(t *testing.T) {
	// Estabelecimento: 09:00–18:00
	// Serviço: 60 minutos
	// O último slot válido deve começar às 17:00 (17:00 + 60min = 18:00).
	// Um slot que começasse às 17:30 terminaria às 18:30 — fora da janela.

	date := makeDate(2026, time.April, 6)
	from := date
	to := date.Add(24 * time.Hour)

	dow := int(date.In(mustLoadLocation(testTimezone)).Weekday())

	biz := bizHoursOpen(dow, "09:00:00", "18:00:00")
	prof := profHoursAvailable(dow, "09:00:00", "18:00:00")

	slots := CalculateSlots(
		biz, prof,
		nil, nil, nil,
		60, from, to,
		testTimezone, "prof-1",
	)

	loc := mustLoadLocation(testTimezone)

	// Último slot deve começar às 17:00.
	if len(slots) == 0 {
		t.Fatal("esperava slots disponíveis, got nenhum")
	}

	lastSlot := slots[len(slots)-1]
	lastLocal := lastSlot.StartsAt.In(loc)
	if lastLocal.Hour() != 17 || lastLocal.Minute() != 0 {
		t.Errorf("último slot deveria começar às 17:00, got %02d:%02d", lastLocal.Hour(), lastLocal.Minute())
	}

	// Nenhum slot deve terminar depois das 18:00.
	closeAt := makeDateTime(2026, time.April, 6, 18, 0)
	for _, s := range slots {
		if s.EndsAt.After(closeAt) {
			t.Errorf("slot %v–%v ultrapassa o fechamento (18:00)", s.StartsAt.In(loc), s.EndsAt.In(loc))
		}
	}

	// 9h às 18h = 9h = 540 min / 60 min = 9 slots
	if len(slots) != 9 {
		t.Errorf("esperava 9 slots, got %d", len(slots))
	}
}

// ─── Cenário 5: dia fechado pelo estabelecimento ─────────────────────────────

func TestCalculateSlots_DiaClosed(t *testing.T) {
	date := makeDate(2026, time.April, 6)
	from := date
	to := date.Add(24 * time.Hour)

	dow := int(date.In(mustLoadLocation(testTimezone)).Weekday())

	// Dia marcado como fechado.
	biz := []BusinessHourRow{
		{DayOfWeek: dow, OpenTime: "09:00:00", CloseTime: "18:00:00", IsClosed: true},
	}
	prof := profHoursAvailable(dow, "09:00:00", "18:00:00")

	slots := CalculateSlots(
		biz, prof,
		nil, nil, nil,
		30, from, to,
		testTimezone, "prof-1",
	)

	if len(slots) != 0 {
		t.Errorf("dia fechado não deve ter slots, got %d", len(slots))
	}
}

// ─── Cenário 6: profissional indisponível ────────────────────────────────────

func TestCalculateSlots_ProfissionalIndisponivel(t *testing.T) {
	date := makeDate(2026, time.April, 6)
	from := date
	to := date.Add(24 * time.Hour)

	dow := int(date.In(mustLoadLocation(testTimezone)).Weekday())

	biz := bizHoursOpen(dow, "09:00:00", "18:00:00")

	// Profissional indisponível no dia.
	prof := []ProfessionalHourRow{
		{DayOfWeek: dow, StartTime: "09:00:00", EndTime: "18:00:00", IsUnavailable: true},
	}

	slots := CalculateSlots(
		biz, prof,
		nil, nil, nil,
		30, from, to,
		testTimezone, "prof-1",
	)

	if len(slots) != 0 {
		t.Errorf("profissional indisponível não deve ter slots, got %d", len(slots))
	}
}

// ─── Cenário 7: bloqueio manual ──────────────────────────────────────────────

func TestCalculateSlots_BloqueioManual(t *testing.T) {
	date := makeDate(2026, time.April, 6)
	from := date
	to := date.Add(24 * time.Hour)

	dow := int(date.In(mustLoadLocation(testTimezone)).Weekday())

	biz := bizHoursOpen(dow, "09:00:00", "18:00:00")
	prof := profHoursAvailable(dow, "09:00:00", "18:00:00")

	// Bloqueio manual: 10:00–12:00
	blockStart := makeDateTime(2026, time.April, 6, 10, 0)
	blockEnd := makeDateTime(2026, time.April, 6, 12, 0)
	blocked := []BlockedPeriod{
		{
			ID:             "block-1",
			ProfessionalID: "prof-1",
			StartsAt:       blockStart,
			EndsAt:         blockEnd,
		},
	}

	slots := CalculateSlots(
		biz, prof,
		nil, blocked, nil,
		30, from, to,
		testTimezone, "prof-1",
	)

	loc := mustLoadLocation(testTimezone)

	// Nenhum slot deve colidir com o bloqueio 10:00–12:00.
	for _, s := range slots {
		if overlaps(s.StartsAt, s.EndsAt, blockStart, blockEnd) {
			t.Errorf("slot %v–%v colide com o bloqueio manual",
				s.StartsAt.In(loc), s.EndsAt.In(loc))
		}
	}

	// 9h às 18h = 18 slots de 30 min.
	// Bloqueio de 10:00–12:00 = 2h = 4 slots bloqueados.
	// Esperado: 18 - 4 = 14 slots.
	if len(slots) != 14 {
		t.Errorf("esperava 14 slots, got %d", len(slots))
	}
}

// ─── Auxiliar ─────────────────────────────────────────────────────────────────

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}
