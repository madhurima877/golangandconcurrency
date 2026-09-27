package main

import "log"

type PaymentProcessor interface {
	Process(amount int) error
}
type RefundProcessor interface {
	Refund(amount int) error
}
type CreditCardProcessor struct {
	balance int
}
type UPIProcessor struct {
	balance int
}

func (p *CreditCardProcessor) Process(amount int) error {
	p.balance += amount
	log.Println(p.balance)
	return nil

}
func (p *CreditCardProcessor) Refund(amount int) error {
	p.balance -= amount
	return nil
}

func (p *UPIProcessor) Process(amount int) error {
	p.balance += amount
	log.Println(p.balance)
	return nil
}
func Pay(p PaymentProcessor, amount int) {
	p.Process(amount)

}
func Refund(r RefundProcessor, amount int) error {
	return r.Refund(amount)
}
func main() {
	processor := &CreditCardProcessor{balance: 0}

	Pay(processor, 100)
	Refund(processor, 20)
	Pay(processor, 10)
}
