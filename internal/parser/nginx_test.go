package parser

import "testing"

type TestCase struct {
	name string
	s    string
	err  bool
}

func TestSplitLine(t *testing.T) {
	t.Helper()
	t.Parallel()

	tests := []TestCase{
		{
			name: "эталонная строка",
			s:    "93.180.71.3 - - [17/May/2015:08:05:32 +0000] \"GET /downloads/product_1 HTTP/1.1\" 304 0 \"-\" \"Debian APT-HTTP/1.3\"",
			err:  false,
		},
		{
			name: "корректный вариант",
			s:    "192.168.1.10 - andrey [05/Nov/2025:23:59:59 +0300] \"POST /api/orders?source=web HTTP/1.1\" 201 1250 \"https://example.com/cart\" \"Mozilla/5.0\"",
			err:  false,
		},
		{
			name: "Нет квадратных скобок вокруг даты",
			s:    "93.180.71.3 - - 17/May/2015:08:05:32 +0000 \"GET /downloads/product_1 HTTP/1.1\" 304 0 \"-\" \"Debian APT-HTTP/1.3\"",
			err:  true},
		{
			name: "Нет кавычек вокруг запроса",
			s:    "93.180.71.3 - - [17/May/2015:08:05:32 +0000] GET /downloads/product_1 HTTP/1.1 304 0 \"-\" \"Debian APT-HTTP/1.3\"",
			err:  true,
		},
		{
			name: "Нет протокола в запросе",
			s:    "93.180.71.3 - - [17/May/2015:08:05:32 +0000] \"GET /downloads/product_1\" 304 0 \"-\" \"Debian APT-HTTP/1.3\"",
			err:  true,
		},
		{
			name: "Нет User-Agent",
			s:    "93.180.71.3 - - [17/May/2015:08:05:32 +0000] \"GET /downloads/product_1 HTTP/1.1\" 304 0 \"-\"",
			err:  true,
		},
		{
			name: "Не закрыта последняя кавычка",
			s:    "93.180.71.3 - - [17/May/2015:08:05:32 +0000] \"GET /downloads/product_1 HTTP/1.1\" 304 0 \"-\" \"Debian APT-HTTP/1.3",
			err:  true,
		},
		{
			name: "Лишнее поле после завершённой записи",
			s:    "93.180.71.3 - - [17/May/2015:08:05:32 +0000] \"GET /downloads/product_1 HTTP/1.1\" 304 0 \"-\" \"Debian APT-HTTP/1.3\" EXTRA",
			err:  true,
		},
		{
			name: "Несуществующий день",
			s:    "93.180.71.3 - - [31/Feb/2015:08:05:32 +0000] \"GET /downloads/product_1 HTTP/1.1\" 304 0 \"-\" \"Debian APT-HTTP/1.3\"",
			err:  true,
		},
		{
			name: "Несуществующий месяц",
			s:    "93.180.71.3 - - [17/Abc/2015:08:05:32 +0000] \"GET /downloads/product_1 HTTP/1.1\" 304 0 \"-\" \"Debian APT-HTTP/1.3\"",
			err:  true,
		},
		{
			name: "Некорректное время",
			s:    "93.180.71.3 - - [17/May/2015:25:05:32 +0000] \"GET /downloads/product_1 HTTP/1.1\" 304 0 \"-\" \"Debian APT-HTTP/1.3\"",
			err:  true,
		},
		{
			name: "Статус содержит буквы",
			s:    "93.180.71.3 - - [17/May/2015:08:05:32 +0000] \"GET /downloads/product_1 HTTP/1.1\" abc 0 \"-\" \"Debian APT-HTTP/1.3\"",
			err:  true,
		},
		{
			name: "Размер содержит буквы",
			s:    "93.180.71.3 - - [17/May/2015:08:05:32 +0000] \"GET /downloads/product_1 HTTP/1.1\" 304 abc \"-\" \"Debian APT-HTTP/1.3\"",
			err:  true,
		},
		{
			name: "Отрицательный размер",
			s:    "93.180.71.3 - - [17/May/2015:08:05:32 +0000] \"GET /downloads/product_1 HTTP/1.1\" 304 -100 \"-\" \"Debian APT-HTTP/1.3\"",
			err:  true,
		},
		{
			name: "Размер больше максимального int64",
			s:    "93.180.71.3 - - [17/May/2015:08:05:32 +0000] \"GET /downloads/product_1 HTTP/1.1\" 304 9223372036854775808 \"-\" \"Debian APT-HTTP/1.3\"",
			err:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := SplitLine(test.s)

			if test.err && err == nil {
				t.Fatal("ожидалась ошибка, но строка определена как корректная")
			}

			if !test.err && err != nil {
				t.Fatalf("строка корректная, но получили ошибку: %v", err)
			}
		})
	}
}
