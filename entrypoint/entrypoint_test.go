package entrypoint

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cli "github.com/stackanetes/kubernetes-entrypoint/client"
	"github.com/stackanetes/kubernetes-entrypoint/logger"
	"github.com/stackanetes/kubernetes-entrypoint/mocks"
)

const (
	testNamespace     = "test"
	dummyResolverName = "dummy"
	loggerInfoText    = "Entrypoint INFO: "
)

var testEntrypoint EntrypointInterface
var testClient cli.ClientInterface

type dummyResolver struct {
	name      string
	namespace string
}

func (d dummyResolver) IsResolved(entry EntrypointInterface) (bool, error) {
	return true, nil
}
func (d dummyResolver) GetName() (name string) {
	return d.name
}

func (d dummyResolver) String() string {
	return fmt.Sprintf("Dummy %s in namespace %s", d.name, d.namespace)
}

func init() {
	testClient = mocks.NewClient()
	testEntrypoint = mocks.NewEntrypointInNamespace(testNamespace)
}

func registerNilResolver() {
	Register(nil)
}

var _ = Describe("Entrypoint", func() {

	dummy := dummyResolver{name: dummyResolverName}

	BeforeEach(func() {
		logger.InfoLog.SetFlags(0)
		logger.WarningLog.SetFlags(0)
		logger.ErrorLog.SetFlags(0)
		logger.OutputJSON = false
	})

	AfterEach(func() {
		// Clear dependencies
		dependencies = make([]Resolver, 0)
	})

	It("registers new nil resolver", func() {
		defer GinkgoRecover()

		Ω(registerNilResolver).Should(Panic())
	})

	It("registers new non-nil resolver", func() {
		defer GinkgoRecover()
		Register(dummy)
		Expect(len(dependencies)).To(Equal(1))
	})

	It("checks Client() method", func() {
		client := testEntrypoint.Client()
		Expect(client).To(Equal(testClient))
	})

	It("resolves main entrypoint with a dummy dependency", func() {
		defer GinkgoRecover()

		// Set output logger to our reader
		r, w, _ := os.Pipe()
		tmp := os.Stdout
		defer func() {
			os.Stdout = tmp
		}()

		logger.InfoLog.SetOutput(w)

		os.Stdout = w
		go func() {
			mainEntrypoint := Entrypoint{client: mocks.NewClient(), namespace: "main"}
			Register(dummy)
			mainEntrypoint.Resolve()
			w.Close()
		}()

		// Wait for resolver to finish
		time.Sleep(5 * time.Second)

		stdout, _ := ioutil.ReadAll(r)
		Expect(string(stdout)).To(Equal(fmt.Sprintf("%sResolving %v\n%sDependency %v is resolved.\n", loggerInfoText, dummy, loggerInfoText, dummy)))
	})
})
