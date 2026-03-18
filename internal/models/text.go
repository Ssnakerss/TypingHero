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
