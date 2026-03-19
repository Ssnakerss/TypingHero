package models

// TextPools содержит коллекцию текстов для каждого уровня сложности (1-10)
// Каждый уровень сложности имеет свой набор текстов, увеличивающихся по длине и сложности
// Используется для генерации случайных текстов в зависимости от выбранного пользователем уровня
var TextPools = map[int][]string{
	1: {
		"The cat sat on the mat.",
		"Dog runs in the park.",
		"Birds fly in the sky.",
		"The sun is bright today.",
		"I like to read books.",
	},
	2: {
		"The quick brown fox jumps over the lazy dog.",
		"She sells seashells by the seashore.",
		"A journey of a thousand miles begins with a single step.",
		"The early bird catches the worm.",
		"Time and tide wait for no man.",
		"How much wood would a woodchuck chuck?",
	},
	3: {
		"Programming is the art of telling a computer what to do.",
		"The best way to predict the future is to create it.",
		"Success is not final, failure is not fatal.",
		"In the middle of difficulty lies opportunity.",
		"Knowledge is power but enthusiasm pulls the switch.",
		"The weather today is pleasant and sunny.",
		"Walking in the park is very relaxing.",
		"Learning a new language takes time and practice.",
	},
	4: {
		"Learning to code is like learning a new language. At first it seems impossible, then it becomes challenging, and finally it becomes natural.",
		"The only way to do great work is to love what you do. If you haven't found it yet, keep looking. Don't settle.",
		"Technology is best when it brings people together. It enables us to connect, share, and learn from each other.",
		"Every expert was once a beginner. Practice makes progress, and consistency is the key to mastery.",
		"Programming is both an art and a science that requires logical thinking.",
		"The ancient castle stood on a hill overlooking the peaceful village below.",
	},
	5: {
		"Software development is not just about writing code; it's about solving problems and creating elegant solutions that make complex systems manageable.",
		"The debugger is twice as hard as writing the code in the first place. Therefore, if you write the code as cleverly as possible, you are by definition not smart enough to debug it.",
		"Any fool can write code that a computer can understand. Good programmers write code that humans can understand. Simplicity is the soul of efficiency.",
		"The scientist carefully conducted the experiment to test her groundbreaking hypothesis.",
	},
	6: {
		"In the world of software, the best code is no code at all. Every new line of code you willingly bring into the world is code that has to be debugged, code that has to be read and understood, and code that has to be supported.",
		"The function of good software is to make the complex appear to be simple to the user. This requires deep understanding of both technology and human psychology.",
		"Programming isn't about what you know; it's about what you can figure out. The only way to go fast, is to go well. Quality is not an act, it is a habit.",
	},
	7: {
		"Design patterns are reusable solutions to commonly occurring problems in software design. They represent best practices evolved over time and provide a standard terminology that makes communication between developers more efficient.",
		"Clean code is happy code. It is writing code that is easy to understand, easy to modify, and easy to extend. The cost of cleaning up code is always less than the cost of maintaining messy code.",
		"Unit testing is not about finding bugs, it is about regression testing. It ensures that changes you make today don't break functionality that worked yesterday.",
	},
	8: {
		"Object-oriented programming was supposed to unify the perspectives of the programmer and the end user. However, modern OOP has become so complex that it often creates more problems than it solves.",
		"The premature optimization is the root of all evil. Yet we should not miss our opportunities to optimize critical sections of code that are executed millions of times.",
		"Dependency injection and inversion of control are powerful patterns that promote loose coupling and make systems more testable and maintainable over time.",
	},
	9: {
		"Functional programming concepts like immutability, higher-order functions, and pure functions can dramatically improve code quality by reducing side effects and making behavior more predictable.",
		"Microservices architecture enables teams to deploy independently, scale horizontally, and adopt different technologies. However, it introduces complexity in distributed systems management.",
		"Event-driven architectures allow systems to be more responsive and loosely coupled. By processing events asynchronously, applications can handle high loads while maintaining responsiveness.",
	},
	10: {
		"Concurrent programming in Go leverages goroutines and channels to create highly efficient and scalable systems. The select statement enables multiplexing between channel operations, while the sync package provides primitives for synchronization.",
		"Distributed systems must handle network partitions, partial failures, and eventual consistency. Understanding CAP theorem trade-offs and implementing proper retry mechanisms with exponential backoff is essential.",
		"Type-driven development and dependent types allow us to encode business rules at the type level, making illegal states unrepresentable and eliminating entire classes of runtime errors through compile-time verification.",
	},
}

// TextPoolsRus содержит коллекцию текстов на русском языке для каждого уровня сложности (1-10)
// Уровни соответствуют возрастающей длине и сложности текстов
// Используется для тренировок на русском языке
var TextPoolsRus = map[int][]string{
	1: {
		"Кот сидит на коврике.",
		"Собака бежит по парку.",
		"Птицы летают в небе.",
		"Сегодня светит солнце.",
		"Я люблю читать книги.",
	},
	2: {
		"Быстрая коричневая лиса прыгает через ленивую собаку.",
		"Шла Саша по шоссе и сосала сушку.",
		"Дорогу осилит идущий.",
		"Утро вечера мудренее.",
		"Век живи — век учись.",
		"Не имей сто рублей, а имей сто друзей.",
	},
	3: {
		"Программирование — это искусство объяснять компьютеру, что нужно делать.",
		"Лучший способ предсказать будущее — создать его.",
		"Успех не случаен, неудача не фатальна.",
		"В середине трудностей скрывается возможность.",
		"Знание — сила, но энтузиазм включает выключатель.",
		"Погода сегодня приятная и солнечная.",
		"Прогулка в парке очень расслабляет.",
		"Изучение нового языка требует времени и практики.",
	},
	4: {
		"Учиться программировать — как учить новый язык. Сначала кажется невозможным, потом сложным, а потом естественным.",
		"Единственный способ сделать великое дело — любить то, что ты делаешь.",
		"Технологии хороши, когда они объединяют людей. Они позволяют общаться, делиться и учиться друг у друга.",
		"Каждый эксперт когда-то был новичком. Практика ведёт к прогрессу, а последовательность — к мастерству.",
		"Программирование — это и искусство, и наука, требующие логического мышления.",
		"Древний замок стоял на холме, откуда открывался вид на тихую деревню.",
	},
	5: {
		"Разработка ПО — это не только написание кода, но и решение задач, создание элегантных решений, упрощающих сложные системы.",
		"Отладка вдвое сложнее, чем написание кода. Поэтому, если вы пишете максимально хитроумный код, вы не сможете его отладить.",
		"Любой дурак может написать код, понятный компьютеру. Хорошие программисты пишут код, понятный людям. Простота — душа эффективности.",
		"Учёный тщательно провёл эксперимент, чтобы проверить свою революционную гипотезу.",
	},
	6: {
		"В мире программирования лучший код — это отсутствие кода. Каждая новая строка — это код, который нужно отлаживать, читать и поддерживать.",
		"Функция хорошего ПО — сделать сложное простым для пользователя. Это требует глубокого понимания технологий и психологии человека.",
		"Программирование — не про знания, а про способность разобраться. Единственный способ двигаться быстро — двигаться хорошо. Качество — не действие, а привычка.",
	},
	7: {
		"Паттерны проектирования — это готовые решения типичных задач. Они отражают лучшие практики и упрощают общение между разработчиками.",
		"Чистый код — это счастливый код. Он легко читается, изменяется и расширяется. Цена за уборку всегда меньше, чем цена за поддержку грязного кода.",
		"Модульное тестирование — не про поиск багов, а про защиту от регрессии. Оно гарантирует, что новые изменения не сломают старую функциональность.",
	},
	8: {
		"Объектно-ориентированное программирование должно было объединить взгляды программиста и пользователя. Но современное ООП стало настолько сложным, что часто создаёт больше проблем, чем решает.",
		"Преждевременная оптимизация — корень всех зол. Однако нельзя упускать возможности оптимизировать критические участки кода, выполняемые миллионы раз.",
		"Внедрение зависимостей и инверсия управления — мощные паттерны, способствующие слабой связанности и улучшающие тестируемость и поддерживаемость систем.",
	},
	9: {
		"Функциональное программирование, с его неизменяемостью, функциями высшего порядка и чистыми функциями, повышает качество кода, снижая побочные эффекты и делая поведение предсказуемым.",
		"Архитектура микросервисов позволяет командам разворачивать независимо, масштабироваться горизонтально и использовать разные технологии. Но она усложняет управление распределёнными системами.",
		"Событийно-ориентированная архитектура делает системы отзывчивее и слабо связанными. Асинхронная обработка событий позволяет выдерживать высокие нагрузки.",
	},
	10: {
		"Параллельное программирование в Go использует горутины и каналы для создания эффективных и масштабируемых систем. Оператор select позволяет мультиплексировать операции с каналами, а пакет sync — обеспечивает синхронизацию.",
		"Распределённые системы должны уметь работать при сетевых разделениях, частичных сбоях и обеспечивать согласованность. Понимание компромиссов CAP-теоремы и реализация механизмов повторных попыток с экспоненциальной задержкой — ключ к надёжности.",
		"Разработка, ориентированная на типы, и зависимые типы позволяют закодировать бизнес-правила на уровне типов, делая недопустимые состояния невозможными и устраняя целые классы ошибок на этапе компиляции.",
	},
}

// GetText возвращает случайный текст для указанного уровня сложности
// Принимает уровень сложности от 1 до 10 и возвращает случайный текст из соответствующего пула
func GetText(difficulty int) string {
	if difficulty < 1 || difficulty > 10 {
		difficulty = 1
	}

	pool := TextPools[difficulty]
	if len(pool) == 0 {
		// Возвращаем текст по умолчанию, если пул пуст
		return "The quick brown fox jumps over the lazy dog."
	}

	// В реальной реализации здесь должна быть случайная выборка
	// Для простоты возвращаем первый текст из пула
	return pool[0]
}
