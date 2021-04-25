package main

import (
	"context"
)

// GolangImp servant implementation
type GolangImp struct {
}

// Init servant init
func (imp *GolangImp) Init() error {
	//initialize servant here:
	//...
	return nil
}

// Destroy servant destory
func (imp *GolangImp) Destroy() {
	//destroy servant here:
	//...
}

func (imp *GolangImp) Ping(ctx context.Context, req string) (string, error) {
	//Doing something in your function
	//...
	return req, nil
}
