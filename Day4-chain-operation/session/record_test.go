package session

import "testing"

var (
	user1 = &User{"Tom", 18}
	user2 = &User{"Sam", 25}
	user3 = &User{"Jack", 25}
)

func testRecordInit(t *testing.T) *Session {
	t.Helper()
	s := NewSession().Model(&User{})
	_ = s.DropTable() // We can ignore the error here, as the table might not exist on the first run.
	err := s.CreateTable()
	if err != nil {
		t.Fatalf("failed on CreateTable: %v", err)
	}
	_, err = s.Insert(user1, user2)
	if err != nil {
		t.Fatalf("failed on Insert: %v", err)
	}
	return s
}

func TestSession_Insert(t *testing.T) {
	s := testRecordInit(t)
	affected, err := s.Insert(user3)
	if err != nil || affected != 1 {
		t.Fatal("failed to create record")
	}
}

func TestSession_Find(t *testing.T) {
	s := testRecordInit(t)
	var users []User
	if err := s.Find(&users); err != nil || len(users) != 2 {
		t.Fatal("failed to query all")
	}
}

func TestSession_First(t *testing.T) {
	s := testRecordInit(t)
	var user User
	if err := s.First(&user); err != nil {
		t.Fatalf("failed to query first: %v", err)
	}
	if user.Name != "Tom" || user.Age != 18 {
		t.Fatalf("failed to get correct record, expected Tom, got %v", user)
	}
}
