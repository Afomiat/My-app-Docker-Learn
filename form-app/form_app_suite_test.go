package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFormApp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FormApp Suite")
}
