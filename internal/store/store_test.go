package store

import (
	"sync"
	"testing"
)

func TestCreateAndGet(t *testing.T) {
	s := New()
	f := Flag{Key: "k", Description: "d", Enabled: true, RolloutPercent: 50}
	s.Create(f)
	got, ok := s.Get("k")
	if !ok {
		t.Fatal("expected flag to exist")
	}
	if got.Key != "k" || got.Description != "d" || !got.Enabled || got.RolloutPercent != 50 {
		t.Fatalf("unexpected flag: %+v", got)
	}
}

func TestCreateUpsert(t *testing.T) {
	s := New()
	s.Create(Flag{Key: "k", Enabled: false})
	s.Create(Flag{Key: "k", Enabled: true, RolloutPercent: 100})
	got, ok := s.Get("k")
	if !ok {
		t.Fatal("expected flag to exist")
	}
	if !got.Enabled || got.RolloutPercent != 100 {
		t.Fatalf("expected upsert to overwrite, got %+v", got)
	}
}

func TestGetMissing(t *testing.T) {
	s := New()
	if _, ok := s.Get("missing"); ok {
		t.Fatal("expected missing flag")
	}
}

func TestListNeverNil(t *testing.T) {
	s := New()
	if s.List() == nil {
		t.Fatal("List must not be nil when empty")
	}
	s.Create(Flag{Key: "a"})
	s.Create(Flag{Key: "b"})
	if got := s.List(); len(got) != 2 {
		t.Fatalf("expected 2 flags, got %d", len(got))
	}
}

func TestUpdate(t *testing.T) {
	s := New()
	s.Create(Flag{Key: "k", Enabled: false})
	got, ok := s.Update("k", Flag{Enabled: true, RolloutPercent: 100})
	if !ok {
		t.Fatal("expected update to succeed")
	}
	if !got.Enabled || got.RolloutPercent != 100 {
		t.Fatalf("unexpected flag: %+v", got)
	}
	if got.Key != "k" {
		t.Fatalf("expected key preserved, got %q", got.Key)
	}
}

func TestUpdateMissing(t *testing.T) {
	s := New()
	if _, ok := s.Update("missing", Flag{Enabled: true}); ok {
		t.Fatal("expected update of missing flag to fail")
	}
}

func TestDelete(t *testing.T) {
	s := New()
	s.Create(Flag{Key: "k"})
	if !s.Delete("k") {
		t.Fatal("expected delete to succeed")
	}
	if _, ok := s.Get("k"); ok {
		t.Fatal("expected flag removed")
	}
}

func TestDeleteMissing(t *testing.T) {
	s := New()
	if s.Delete("missing") {
		t.Fatal("expected delete of missing flag to fail")
	}
}

func TestConcurrentAccess(t *testing.T) {
	s := New()
	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := string(rune('a' + i%26))
			s.Create(Flag{Key: key, Enabled: true, RolloutPercent: i % 101})
			s.Get(key)
			s.List()
			s.Update(key, Flag{Enabled: false})
			s.Delete(key)
		}(i)
	}
	wg.Wait()
}
