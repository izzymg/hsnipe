package web

import "time"

/** Web: computerlounge.co.nz */

// const domain = "https://www.computerlounge.co.nz"
const httpTimeout = 10 * time.Second

func CreateComputerLoungeProvider() *ComputerLoungeProvider {
	return &ComputerLoungeProvider{
		client: createClient(10*time.Second, "https://www.computerlounge.co.nz"),
	}
}

type ComputerLoungeProvider struct {
	client *webClient
}

func (c ComputerLoungeProvider) Name() string {
	return "Computer Lounge"
}

func (c ComputerLoungeProvider) SearchPage(query string, page int) ([]Product, error) {
	return nil, nil
}
