package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
)

func main() {
	/*exKernel := excel_analysis.NewExccelAnalysis("./config/client_secret.json", "1MjTmlm-4Inf4u0kluKfuCcrSRId1CPbLVHTLGr3lFRU")
	data, _ := exKernel.GetHeaderInfo()
	dataStr, _ := json.Marshal(data)
	fmt.Println(string(dataStr))*/

	file, err := excelize.OpenFile("ServerTable.xlsx")
	if err != nil {
		logrus.Fatal(err.Error())
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			{
				fmt.Println(err)
			}
		}
	}()

	rows, err := file.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
	}

	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}

		fmt.Println()
	}

	if err := initConfig(); err != nil {
		logrus.Fatalf("Ошибка при инициализации файла конфигурации: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Ошибка при инициализации глобальных переменных: %s", err.Error())
	}

	// Инициализация логгера
	logrus.SetFormatter(new(logrus.JSONFormatter))

	fileError, err := os.OpenFile(viper.GetString("paths.logs.error"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logrus.AddHook(&writer.Hook{
			Writer: fileError,
			LogLevels: []logrus.Level{
				logrus.ErrorLevel,
			},
		})
	} else {
		logrus.SetOutput(os.Stderr)
		logrus.Error("Failed to log to file, using default stderr")
	}

	defer fileError.Close()

	fileInfo, err := os.OpenFile(viper.GetString("paths.logs.info"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logrus.AddHook(&writer.Hook{
			Writer: fileInfo,
			LogLevels: []logrus.Level{
				logrus.InfoLevel,
				logrus.DebugLevel,
			},
		})
	} else {
		logrus.SetOutput(os.Stderr)
		logrus.Error("Failed to log to file, using default stderr")
	}

	defer fileInfo.Close()

	fileWarn, err := os.OpenFile(viper.GetString("paths.logs.warn"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logrus.AddHook(&writer.Hook{
			Writer: fileWarn,
			LogLevels: []logrus.Level{
				logrus.WarnLevel,
			},
		})
	} else {
		logrus.SetOutput(os.Stderr)
		logrus.Error("Failed to log to file, using default stderr")
	}

	defer fileWarn.Close()

	fileFatal, err := os.OpenFile(viper.GetString("paths.logs.fatal"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logrus.AddHook(&writer.Hook{
			Writer: fileFatal,
			LogLevels: []logrus.Level{
				logrus.FatalLevel,
			},
		})
	} else {
		logrus.SetOutput(os.Stderr)
		logrus.Error("Failed to log to file, using default stderr")
	}

	defer fileFatal.Close()
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
