package utilities

import (
	"fmt"
	"math/rand"
	"time"

	guuid "github.com/google/uuid"
)

func RandomUUID() string {
	return guuid.NewString()
}

func GenerateRefId() (refId string) {
	rand.Seed(time.Now().UTC().UnixNano())
	generateCode := 100000 + rand.Intn(999999-100000)
	refId = "PPOB" + fmt.Sprint(generateCode)
	return refId
}

func GenerateEmail() (email string) {
	rand.Seed(time.Now().UTC().UnixNano())
	generateCode := 100000 + rand.Intn(999999-100000)
	email = "test" + fmt.Sprint(generateCode) + "@gmail.com"
	return email
}

func GenerateNoTagihan() (refId string) {
	rand.Seed(time.Now().UTC().UnixNano())
	generateCode := 100000 + rand.Intn(999999-100000)
	refId = "TAGIHAN/" + fmt.Sprint(generateCode)
	return refId
}
