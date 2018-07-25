package rpc

import "testing"

func TestMailbox(t *testing.T) {
	mb := NewMailbox(1, 0, 1)
	if mb.ServiceId() != 1 || mb.Flag() != 0 || mb.Id() != 1 {
		t.Fatalf("mailbox error, want: 1, 0, 1, have:%d, %d, %d", mb.ServiceId(), mb.Flag(), mb.Id())
	}

	mb1 := NewMailbox(1, 1, 1000)
	if mb1.Flag() != 1 || !mb1.IsClient() || mb1.IsObject() {
		t.Fatalf("mailbox error, want: 1, true, false, have:%d, %v, %v", mb1.Flag(), mb1.IsClient(), mb1.IsObject())
	}

	mb2 := mb.NewObjectId(2, 3, 4)
	if mb2.ServiceId() != 1 || mb2.Flag() != 0 || !mb2.IsObject() || mb2.IsClient() {
		t.Fatalf("mailbox error, want: 1, 0, true, false, have:%d, %d, %v, %v", mb2.ServiceId(), mb2.Flag(), mb2.IsObject(), mb2.IsClient())
	}

	if mb2.Identity() != 2 || mb2.ObjectIndex() != 4 {
		t.Fatalf("mailbox error, want: 2, 4, have:%d, %d", mb2.Identity(), mb2.ObjectIndex())
	}
}
