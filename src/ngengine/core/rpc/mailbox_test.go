package rpc

import "testing"

func TestMailbox(t *testing.T) {
	mb := NewMailbox(1, 1, 1)
	if mb.ServiceId() != 1 || mb.Flag() != 1 || mb.Id() != 1 {
		t.Fatalf("mailbox error, want: 1, 1, 1, have:%d, %d, %d", mb.ServiceId(), mb.Flag(), mb.Id())
	}

	mb2 := mb.NewObjectId(2, 3, 4)
	if mb2.ServiceId() != 1 || mb2.Flag() != 1 {
		t.Fatalf("mailbox error, want: 1, 1, have:%d, %d", mb2.ServiceId(), mb2.Flag())
	}

	if mb2.ObjectType() != 2 || mb2.ObjectIndex() != 4 {
		t.Fatalf("mailbox error, want: 2, 4, have:%d, %d", mb2.ObjectType(), mb2.ObjectIndex())
	}
}
