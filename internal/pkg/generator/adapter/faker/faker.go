package faker

import (
	"github.com/bxcodec/faker/v3"
)

func FakeData(a interface{}) error {
	return faker.FakeData(a)
}
