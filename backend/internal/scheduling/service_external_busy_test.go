package scheduling

import (
	"context"
	"testing"
	"time"

	"schedule/internal/shared"
)

type fakeSchedulingRepo struct{}

func (fakeSchedulingRepo) GetBusinessHours(ctx context.Context, establishmentID string) ([]BusinessHourRow, error) {
	return nil, nil
}

func (fakeSchedulingRepo) GetProfessionalHours(ctx context.Context, professionalID string) ([]ProfessionalHourRow, error) {
	return []ProfessionalHourRow{
		{DayOfWeek: 1, StartTime: "09:00:00", EndTime: "18:00:00", IsUnavailable: false},
	}, nil
}

func (fakeSchedulingRepo) GetConfirmedAppointments(ctx context.Context, professionalID string, from, to time.Time) ([]Appointment, error) {
	return nil, nil
}

func (fakeSchedulingRepo) GetBlockedPeriods(ctx context.Context, professionalID string, from, to time.Time) ([]BlockedPeriod, error) {
	return nil, nil
}

func (fakeSchedulingRepo) GetActiveProfessionals(ctx context.Context, establishmentID string, serviceID string) ([]ProfessionalRef, error) {
	return nil, nil
}

func (fakeSchedulingRepo) GetServiceDuration(ctx context.Context, serviceID, establishmentID string) (int, error) {
	return 30, nil
}

func (fakeSchedulingRepo) GetServiceName(ctx context.Context, serviceID, establishmentID string) (string, error) {
	return "Servico teste", nil
}

func (fakeSchedulingRepo) GetEstablishmentTimezone(ctx context.Context, establishmentID string) (string, error) {
	return testTimezone, nil
}

func (fakeSchedulingRepo) GetMinAdvanceCancelHours(ctx context.Context, establishmentID string) (int, error) {
	return 0, nil
}

func (fakeSchedulingRepo) ValidateProfessionalForService(ctx context.Context, establishmentID, professionalID, serviceID string) error {
	return nil
}

func (fakeSchedulingRepo) FindAppointmentByIdempotencyKey(ctx context.Context, establishmentID, idempotencyKey string) (*Appointment, error) {
	return nil, nil
}

func (fakeSchedulingRepo) FindAppointmentByIDAndPhone(ctx context.Context, establishmentID, appointmentID, phone string) (*Appointment, error) {
	return nil, nil
}

func (fakeSchedulingRepo) CreateAppointment(ctx context.Context, input CreateAppointmentInput, endsAt time.Time) (*Appointment, error) {
	return nil, nil
}

func (fakeSchedulingRepo) SetAppointmentGoogleEventID(ctx context.Context, establishmentID, appointmentID string, googleEventID *string) error {
	return nil
}

func (fakeSchedulingRepo) CancelAppointment(ctx context.Context, establishmentID, appointmentID string) (*Appointment, error) {
	return nil, nil
}

func (fakeSchedulingRepo) RescheduleAppointment(ctx context.Context, current *Appointment, input CreateAppointmentInput, endsAt time.Time) (*Appointment, error) {
	return nil, nil
}

func (fakeSchedulingRepo) ListManagerAppointments(ctx context.Context, establishmentID string, filter AppointmentFilter) ([]ManagerAppointmentRow, int, error) {
	return nil, 0, nil
}

func (fakeSchedulingRepo) GetManagerAppointment(ctx context.Context, establishmentID, appointmentID string) (*ManagerAppointmentRow, error) {
	return nil, nil
}

func (fakeSchedulingRepo) UpdateAppointmentStatus(ctx context.Context, input UpdateStatusInput) (*ManagerAppointmentRow, error) {
	return nil, nil
}

func (fakeSchedulingRepo) ListManagerBlockedPeriods(ctx context.Context, establishmentID, professionalID, dateFrom, dateTo string) ([]ManagerBlockedPeriod, error) {
	return nil, nil
}

func (fakeSchedulingRepo) CreateBlockedPeriod(ctx context.Context, input CreateBlockedPeriodInput) (*ManagerBlockedPeriod, error) {
	return nil, nil
}

func (fakeSchedulingRepo) DeleteBlockedPeriod(ctx context.Context, establishmentID, blockedPeriodID string) error {
	return nil
}

type fakeBusyProvider struct {
	called bool
	period []Period
	err    error
}

func (f *fakeBusyProvider) ListBusyPeriods(ctx context.Context, establishmentID, professionalID string, from, to time.Time) ([]Period, error) {
	f.called = true
	return f.period, f.err
}

type fakeAppointmentEventSyncer struct {
	createCalled bool
	deleteCalled bool
	eventID      string
	err          error
}

func (f *fakeAppointmentEventSyncer) ListBusyPeriods(ctx context.Context, establishmentID, professionalID string, from, to time.Time) ([]Period, error) {
	return nil, nil
}

func (f *fakeAppointmentEventSyncer) CreateAppointmentEvent(ctx context.Context, appointment *Appointment) (string, error) {
	f.createCalled = true
	return f.eventID, f.err
}

func (f *fakeAppointmentEventSyncer) DeleteAppointmentEvent(ctx context.Context, appointment *Appointment) error {
	f.deleteCalled = true
	return f.err
}

func TestGetSlotsForProfessional_UsesExternalBusyProvider(t *testing.T) {
	provider := &fakeBusyProvider{
		period: []Period{
			{
				StartsAt: makeDateTime(2026, time.April, 6, 10, 0),
				EndsAt:   makeDateTime(2026, time.April, 6, 11, 0),
			},
		},
	}

	service := NewService(fakeSchedulingRepo{}, nil, provider)
	date := makeDate(2026, time.April, 6)
	dow := int(date.In(mustLoadLocation(testTimezone)).Weekday())

	slots, err := service.getSlotsForProfessional(
		context.Background(),
		"prof-1",
		bizHoursOpen(dow, "09:00:00", "18:00:00"),
		30,
		AvailabilityOptions{
			EstablishmentID: "est-1",
			ProfessionalID:  "prof-1",
			DateFrom:        date,
			DateTo:          date,
			Timezone:        testTimezone,
		},
	)
	if err != nil {
		t.Fatalf("esperava sucesso, got err=%v", err)
	}

	if !provider.called {
		t.Fatal("esperava consulta ao provedor de busy externo")
	}

	for _, slot := range slots {
		if overlaps(slot.StartsAt, slot.EndsAt, provider.period[0].StartsAt, provider.period[0].EndsAt) {
			t.Fatalf("slot %v-%v nao deveria ignorar busy externo", slot.StartsAt, slot.EndsAt)
		}
	}
}

func TestCreatePublicAppointment_IgnoresIntegrationNotConfiguredFromBusyProvider(t *testing.T) {
	startsAt := makeDateTime(2026, time.April, 6, 9, 0)
	provider := &fakeBusyProvider{err: shared.ErrIntegrationNotConfigured}
	repo := &publicAppointmentRepoStub{
		duration: 30,
		timezone: testTimezone,
		created: &Appointment{
			ID:              "appt-1",
			EstablishmentID: "est-1",
			ProfessionalID:  "prof-1",
			ServiceID:       "svc-1",
			ClientName:      "Maria",
			ClientEmail:     strPtr("maria@example.com"),
			ClientPhone:     "11999999999",
			ClientBirthDate: strPtr("1990-01-01"),
			StartsAt:        startsAt,
			EndsAt:          startsAt.Add(30 * time.Minute),
			Status:          "confirmed",
		},
	}

	service := NewService(repo, nil, provider)

	result, err := service.CreatePublicAppointment(context.Background(), CreateAppointmentInput{
		EstablishmentID: "est-1",
		ServiceID:       "svc-1",
		ProfessionalID:  "prof-1",
		StartsAt:        startsAt,
		ClientName:      "Maria",
		ClientEmail:     "maria@example.com",
		ClientPhone:     "11999999999",
		ClientBirthDate: "1990-01-01",
		IdempotencyKey:  "idem-1",
	})
	if err != nil {
		t.Fatalf("esperava sucesso sem integracao configurada, got err=%v", err)
	}

	if result == nil || result.ID != "appt-1" {
		t.Fatalf("esperava agendamento criado, got %#v", result)
	}

	if !provider.called {
		t.Fatal("esperava tentativa de consulta ao provedor de busy externo")
	}
}

func TestCreatePublicAppointment_SyncsGoogleEventIDWhenSyncSucceeds(t *testing.T) {
	startsAt := makeDateTime(2026, time.April, 6, 9, 0)
	repo := &publicAppointmentRepoStub{
		duration: 30,
		timezone: testTimezone,
		created: &Appointment{
			ID:              "appt-2",
			EstablishmentID: "est-1",
			ProfessionalID:  "prof-1",
			ServiceID:       "svc-1",
			ClientName:      "Paula",
			ClientEmail:     strPtr("paula@example.com"),
			ClientPhone:     "11911111111",
			ClientBirthDate: strPtr("1991-01-01"),
			StartsAt:        startsAt,
			EndsAt:          startsAt.Add(30 * time.Minute),
			Status:          "confirmed",
		},
	}
	syncer := &fakeAppointmentEventSyncer{eventID: "gcal-1"}

	service := NewService(repo, nil, syncer)

	result, err := service.CreatePublicAppointment(context.Background(), CreateAppointmentInput{
		EstablishmentID: "est-1",
		ServiceID:       "svc-1",
		ProfessionalID:  "prof-1",
		StartsAt:        startsAt,
		ClientName:      "Paula",
		ClientEmail:     "paula@example.com",
		ClientPhone:     "11911111111",
		ClientBirthDate: "1991-01-01",
		IdempotencyKey:  "idem-2",
	})
	if err != nil {
		t.Fatalf("esperava sucesso com syncer, got err=%v", err)
	}
	if result == nil || result.ID != "appt-2" {
		t.Fatalf("esperava resultado do agendamento criado, got %#v", result)
	}
	if !syncer.createCalled {
		t.Fatal("esperava criacao de evento no Google Calendar")
	}
	if repo.lastGoogleEventID == nil || *repo.lastGoogleEventID != "gcal-1" {
		t.Fatalf("esperava persistir google_event_id, got %#v", repo.lastGoogleEventID)
	}
}

type publicAppointmentRepoStub struct {
	duration          int
	timezone          string
	created           *Appointment
	lastGoogleEventID *string
}

func (s *publicAppointmentRepoStub) GetBusinessHours(ctx context.Context, establishmentID string) ([]BusinessHourRow, error) {
	date := makeDate(2026, time.April, 6)
	dow := int(date.In(mustLoadLocation(testTimezone)).Weekday())
	return bizHoursOpen(dow, "09:00:00", "18:00:00"), nil
}

func (s *publicAppointmentRepoStub) GetProfessionalHours(ctx context.Context, professionalID string) ([]ProfessionalHourRow, error) {
	date := makeDate(2026, time.April, 6)
	dow := int(date.In(mustLoadLocation(testTimezone)).Weekday())
	return []ProfessionalHourRow{
		{DayOfWeek: dow, StartTime: "09:00:00", EndTime: "18:00:00", IsUnavailable: false},
	}, nil
}

func (s *publicAppointmentRepoStub) GetConfirmedAppointments(ctx context.Context, professionalID string, from, to time.Time) ([]Appointment, error) {
	return nil, nil
}

func (s *publicAppointmentRepoStub) GetBlockedPeriods(ctx context.Context, professionalID string, from, to time.Time) ([]BlockedPeriod, error) {
	return nil, nil
}

func (s *publicAppointmentRepoStub) GetActiveProfessionals(ctx context.Context, establishmentID string, serviceID string) ([]ProfessionalRef, error) {
	return nil, nil
}

func (s *publicAppointmentRepoStub) GetServiceDuration(ctx context.Context, serviceID, establishmentID string) (int, error) {
	return s.duration, nil
}

func (s *publicAppointmentRepoStub) GetServiceName(ctx context.Context, serviceID, establishmentID string) (string, error) {
	return "Servico teste", nil
}

func (s *publicAppointmentRepoStub) GetEstablishmentTimezone(ctx context.Context, establishmentID string) (string, error) {
	return s.timezone, nil
}

func (s *publicAppointmentRepoStub) GetMinAdvanceCancelHours(ctx context.Context, establishmentID string) (int, error) {
	return 0, nil
}

func (s *publicAppointmentRepoStub) ValidateProfessionalForService(ctx context.Context, establishmentID, professionalID, serviceID string) error {
	return nil
}

func (s *publicAppointmentRepoStub) FindAppointmentByIdempotencyKey(ctx context.Context, establishmentID, idempotencyKey string) (*Appointment, error) {
	return nil, shared.ErrNotFound
}

func (s *publicAppointmentRepoStub) FindAppointmentByIDAndPhone(ctx context.Context, establishmentID, appointmentID, phone string) (*Appointment, error) {
	return nil, shared.ErrNotFound
}

func (s *publicAppointmentRepoStub) CreateAppointment(ctx context.Context, input CreateAppointmentInput, endsAt time.Time) (*Appointment, error) {
	return s.created, nil
}

func (s *publicAppointmentRepoStub) SetAppointmentGoogleEventID(ctx context.Context, establishmentID, appointmentID string, googleEventID *string) error {
	s.lastGoogleEventID = googleEventID
	return nil
}

func (s *publicAppointmentRepoStub) CancelAppointment(ctx context.Context, establishmentID, appointmentID string) (*Appointment, error) {
	return nil, shared.ErrNotFound
}

func (s *publicAppointmentRepoStub) RescheduleAppointment(ctx context.Context, current *Appointment, input CreateAppointmentInput, endsAt time.Time) (*Appointment, error) {
	return nil, shared.ErrNotFound
}

func (s *publicAppointmentRepoStub) ListManagerAppointments(ctx context.Context, establishmentID string, filter AppointmentFilter) ([]ManagerAppointmentRow, int, error) {
	return nil, 0, nil
}

func (s *publicAppointmentRepoStub) GetManagerAppointment(ctx context.Context, establishmentID, appointmentID string) (*ManagerAppointmentRow, error) {
	return nil, shared.ErrNotFound
}

func (s *publicAppointmentRepoStub) UpdateAppointmentStatus(ctx context.Context, input UpdateStatusInput) (*ManagerAppointmentRow, error) {
	return nil, shared.ErrNotFound
}

func (s *publicAppointmentRepoStub) ListManagerBlockedPeriods(ctx context.Context, establishmentID, professionalID, dateFrom, dateTo string) ([]ManagerBlockedPeriod, error) {
	return nil, nil
}

func (s *publicAppointmentRepoStub) CreateBlockedPeriod(ctx context.Context, input CreateBlockedPeriodInput) (*ManagerBlockedPeriod, error) {
	return nil, nil
}

func (s *publicAppointmentRepoStub) DeleteBlockedPeriod(ctx context.Context, establishmentID, blockedPeriodID string) error {
	return shared.ErrNotFound
}

func strPtr(value string) *string {
	return &value
}
