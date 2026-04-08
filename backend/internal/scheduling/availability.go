package scheduling

// availability.go — Motor de disponibilidade.
// Este arquivo é o coração do produto. Qualquer alteração exige testes unitários.
// Ver internal/scheduling/availability_test.go.
//
// Algoritmo principal: CalculateSlots
// Regras (em ordem de prioridade):
//   1. Horário de funcionamento do estabelecimento
//   2. Jornada individual do profissional no dia
//   3. Duração do serviço deve caber inteiramente no slot
//   4. Agendamentos já confirmados (status != 'cancelled')
//   5. Bloqueios manuais cadastrados pelo gestor
//   6. Eventos externos do Google Agenda (quando integrado)

import (
	"fmt"
	"time"
)

// Period representa um intervalo de tempo ocupado (usado para eventos externos).
type Period struct {
	StartsAt time.Time
	EndsAt   time.Time
}

// overlaps retorna true se dois intervalos se sobrepõem.
// Colisão = starts_at < slotEnd AND ends_at > slotStart
func overlaps(aStart, aEnd, bStart, bEnd time.Time) bool {
	return aStart.Before(bEnd) && aEnd.After(bStart)
}

// parseTimeOfDay converte uma string "HH:MM:SS" ou "HH:MM" no timezone dado
// para um time.Time no dia especificado (em UTC).
func parseTimeOfDay(dayUTC time.Time, timeStr string, loc *time.Location) (time.Time, error) {
	// Constrói uma data no timezone do estabelecimento com o horário informado.
	year, month, day := dayUTC.In(loc).Date()

	var h, m, s int
	switch len(timeStr) {
	case 8: // HH:MM:SS
		_, err := fmt.Sscanf(timeStr, "%d:%d:%d", &h, &m, &s)
		if err != nil {
			return time.Time{}, fmt.Errorf("formato de horário inválido: %q", timeStr)
		}
	case 5: // HH:MM
		_, err := fmt.Sscanf(timeStr, "%d:%d", &h, &m)
		if err != nil {
			return time.Time{}, fmt.Errorf("formato de horário inválido: %q", timeStr)
		}
	default:
		return time.Time{}, fmt.Errorf("formato de horário inválido: %q", timeStr)
	}

	t := time.Date(year, month, day, h, m, s, 0, loc)
	return t.UTC(), nil
}

// CalculateSlots calcula os slots disponíveis para um profissional em um intervalo.
//
// É uma função pura: recebe os dados já carregados, sem acesso ao banco.
//
// Parâmetros:
//   - bizHours: horários de funcionamento do estabelecimento (indexados por DayOfWeek 0-6)
//   - profHours: jornada individual do profissional (indexados por DayOfWeek 0-6)
//   - appointments: agendamentos confirmados no período
//   - blockedPeriods: bloqueios manuais no período
//   - externalBusy: eventos do Google Calendar (pode ser nil)
//   - durationMinutes: duração do serviço em minutos
//   - from, to: intervalo de datas em UTC (from inclusivo, to exclusivo)
//   - timezone: ex: "America/Sao_Paulo"
//   - professionalID: ID do profissional (incluído em cada Slot retornado)
func CalculateSlots(
	bizHours []BusinessHourRow,
	profHours []ProfessionalHourRow,
	appointments []Appointment,
	blockedPeriods []BlockedPeriod,
	externalBusy []Period,
	durationMinutes int,
	from, to time.Time,
	timezone string,
	professionalID string,
) []Slot {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		// Fallback para UTC se timezone inválido
		loc = time.UTC
	}

	duration := time.Duration(durationMinutes) * time.Minute

	// Indexa business_hours e professional_hours por DayOfWeek para O(1) lookup.
	bizByDay := make(map[int]BusinessHourRow, len(bizHours))
	for _, bh := range bizHours {
		bizByDay[bh.DayOfWeek] = bh
	}

	profByDay := make(map[int]ProfessionalHourRow, len(profHours))
	for _, ph := range profHours {
		profByDay[ph.DayOfWeek] = ph
	}

	var slots []Slot

	// Itera dia a dia no intervalo [from, to).
	// Normaliza 'from' para o início do dia em UTC.
	current := truncateToDay(from, loc)
	end := truncateToDay(to, loc)
	if to.After(truncateToDay(to, loc)) {
		// Se 'to' não é exatamente meia-noite, inclui o dia de 'to'.
		end = truncateToDay(to, loc).Add(24 * time.Hour)
	}

	for current.Before(end) {
		dayInLoc := current.In(loc)
		dow := int(dayInLoc.Weekday()) // 0=domingo, 6=sábado

		// Regra 1: Verificar horário de funcionamento do estabelecimento.
		bh, hasBiz := bizByDay[dow]
		if !hasBiz || bh.IsClosed {
			current = current.Add(24 * time.Hour)
			continue
		}

		bizOpen, err := parseTimeOfDay(current, bh.OpenTime, loc)
		if err != nil {
			current = current.Add(24 * time.Hour)
			continue
		}
		bizClose, err := parseTimeOfDay(current, bh.CloseTime, loc)
		if err != nil {
			current = current.Add(24 * time.Hour)
			continue
		}

		// Regra 2: Verificar jornada individual do profissional no dia.
		ph, hasProf := profByDay[dow]
		if hasProf && ph.IsUnavailable {
			current = current.Add(24 * time.Hour)
			continue
		}

		// Janela de trabalho: interseção entre horário do estabelecimento e jornada do profissional.
		windowStart := bizOpen
		windowEnd := bizClose

		if hasProf {
			profStart, err := parseTimeOfDay(current, ph.StartTime, loc)
			if err == nil && profStart.After(windowStart) {
				windowStart = profStart
			}
			profEnd, err := parseTimeOfDay(current, ph.EndTime, loc)
			if err == nil && profEnd.Before(windowEnd) {
				windowEnd = profEnd
			}
		}

		// Janela inválida: profissional começa depois do fechamento ou vice-versa.
		if !windowStart.Before(windowEnd) {
			current = current.Add(24 * time.Hour)
			continue
		}

		// Regra 3: Dividir janela em slots de durationMinutes.
		slotStart := windowStart
		for {
			slotEnd := slotStart.Add(duration)

			// O slot deve caber inteiramente dentro da janela de trabalho.
			if slotEnd.After(windowEnd) {
				break
			}

			// Regra 4: Verificar colisão com agendamentos confirmados.
			busy := false
			for _, appt := range appointments {
				if overlaps(slotStart, slotEnd, appt.StartsAt, appt.EndsAt) {
					busy = true
					break
				}
			}

			// Regra 5: Verificar colisão com bloqueios manuais.
			if !busy {
				for _, bp := range blockedPeriods {
					if overlaps(slotStart, slotEnd, bp.StartsAt, bp.EndsAt) {
						busy = true
						break
					}
				}
			}

			// Regra 6: Verificar colisão com eventos externos (Google Calendar).
			if !busy {
				for _, ext := range externalBusy {
					if overlaps(slotStart, slotEnd, ext.StartsAt, ext.EndsAt) {
						busy = true
						break
					}
				}
			}

			if !busy {
				slots = append(slots, Slot{
					StartsAt:       slotStart,
					EndsAt:         slotEnd,
					ProfessionalID: professionalID,
				})
			}

			slotStart = slotEnd
		}

		current = current.Add(24 * time.Hour)
	}

	return slots
}

// truncateToDay retorna o instante UTC correspondente à meia-noite do dia
// ao qual 't' pertence no timezone 'loc'.
func truncateToDay(t time.Time, loc *time.Location) time.Time {
	y, m, d := t.In(loc).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, loc).UTC()
}
