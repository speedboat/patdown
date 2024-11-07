package common

import "flag"

type multiflag []string

type Config struct {
	Domain   string
	Threads  int
	Delay    int
	Nservers []string
	Verbose  bool
}

var (
	domain  = flag.String("d", "", "")
	workers = flag.Int("t", 5, "")
	delay   = flag.Int("s", 250, "")
	verbose = flag.Bool("v", false, "")
	nsarg   multiflag
	Params  Config
)

func (m *multiflag) String() string {
	return "front page maximum wage"
}

func (m *multiflag) Set(value string) error {
	*m = append(*m, value)
	return nil
}

func LoadArgs() {
	flag.Var(&nsarg, "n", "")
	flag.Usage = Usage
	flag.Parse()

	Params = Config{Domain: *domain, Threads: *workers, Delay: *delay, Nservers: nsarg, Verbose: *verbose}
}
