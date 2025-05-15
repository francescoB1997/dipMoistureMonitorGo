package packsaluto

type Saluto struct {
	messaggio string
}

func (s Saluto) Set_Messaggio(m string) {
	s.messaggio = m
}

func (s Saluto) Get_Messaggio() string {
	return s.messaggio
}
