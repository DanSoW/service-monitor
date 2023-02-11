package main

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

/* Модель для сохранения результатов выполнения команды по SSH*/
type ResultExec struct {
	Command  string
	ExitCode int
	Output   string
}

/* Модель для сохранения множества результатов выполнения команд по SSH */
type ResultExecCmd struct {
	Results []ResultExec
}

func (rec *ResultExecCmd) ToString() string {
	var result string
	result = ""

	for _, item := range rec.Results {
		result += fmt.Sprintf("%s\n", item.Output)
	}

	return result
}

func ExecCommand(client *ssh.Client, command string) (*ResultExec, error) {
	// Создание новой сессии клиента
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	// Выполнение команды в рамках клиентской сессии и возвращение результата в переменную
	output, err := session.Output(command)
	if err != nil {
		return nil, err
	}

	defer session.Close()

	return &ResultExec{
		Command:  command,
		ExitCode: 0,
		Output:   string(output),
	}, nil
}

/* Покдлючение к удалённому серверу по SSH */
func ExecCommands(host string, port int, user string, password string, commands []string) (*ResultExecCmd, error) {
	// Создание конфигурации SSH для подключения к серверу
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Формирование адреса для подключения
	addr := fmt.Sprintf("%s:%d", host, port)

	// Создание подключения по SSH
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	var results []ResultExec

	for _, item := range commands {
		result, err := ExecCommand(client, item)
		if err != nil {
			return nil, err
		}

		results = append(results, *result)
	}

	return &ResultExecCmd{
		Results: results,
	}, nil
}
