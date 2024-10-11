package MongoDB

import (
	"fmt"
	"os"
	"testing"
)

var (
	Uri            = os.Getenv("DB_URI")
	databaseName   = os.Getenv("DB_NAME")
	collectionName = os.Getenv("DB_COLLECTION_NAME")
)

func Test_Health(t *testing.T) {
	srv, err := InitMongo(Uri, databaseName, collectionName)
	if err != nil {
		t.Fatal(err)
	}
	stats := srv.Health()

	if stats["message"] != "It's healthy" {
		t.Fatalf("expected message to be 'It's healthy', got %s", stats["message"])
	}
}

func TestAddGetDelete(t *testing.T) {
	fmt.Println(Uri, databaseName, collectionName)
	db, err := InitMongo(Uri, databaseName, collectionName)
	if err != nil {
		t.Fatal(err)
	}
	err = db.Add("tg-id", "token", "refreshToken")
	if err != nil {
		t.Fatal(err)
	}
	usr, err := db.Get("tg-id")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("user", usr)
	//err = db.Delete("tg-id")
	//if err != nil {
	//	t.Fatal(err)
	//}
}
func TestUpdate(t *testing.T) {
	db, err := InitMongo(Uri, databaseName, collectionName)
	if err != nil {
		t.Fatal(err)
	}
	err = db.Add("tg-id", "token", "refreshToken")
	if err != nil {
		t.Fatal(err)
	}
	usr, err := db.Get("tg-id")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("user1 - ", usr)
	err = db.Update("tg-id", "token1", "refreshToken1")
	if err != nil {
		t.Fatal(err)
	}
	usr, err = db.Get("tg-id")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("user2 - ", usr)
	err = db.Delete("tg-id")
	if err != nil {
		t.Fatal(err)
	}
}
