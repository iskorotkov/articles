# Go: использование значений VS указателей

- [Go: использование значений VS указателей](#go-использование-значений-vs-указателей)
  - [Читаемость](#читаемость)
    - [Параметр](#параметр)
    - [Возвращаемое значение](#возвращаемое-значение)
  - [Производительность](#производительность)
    - [Fizzbuzz](#fizzbuzz)
  - [Расширяемость](#расширяемость)
  - [Консистентность](#консистентность)
  - [Потенциальные ошибки](#потенциальные-ошибки)

Go - один из немногих языков, в которых структуры можно передавать параметрами и возвращать из функций как по значению, так и по указателю. Это приводит к большей выразительности языка, но также разделяет общество разработчиков Go на два лагеря: сторонников указателей и сторонников значений.

В данной статье предлагается во многом субъективное сравнение обоих способов и делается попытка убедить читателей передавать и возвращать значения в тех случаях, где это возможно.

## Читаемость

### Параметр

Самое простое и, наверное, самое главное свойство кода - это его читаемость. Указатели имеют больше возможных значений, так что работать с ними оказывается немного сложнее, чем со значениями. Например:

```go
func foo(ctx context.Context, user *User, order *Order) (*Receipt, error) {
    // ...
}
```

При рефакторинге такой функции будет проскальзывать мысль: а не может ли одним из параметров функции быть передан `nil`? Даже если это маловероятно и разработчики договорились/из документации следует, что `nil` передан не будет, есть ли гарантии, что `nil` не попадет сюда по ошибке? Кто-то из разработчиков может добавить новую функциональность, вызывающую данную функцию, и забыть добавить проверку на `nil`. А может из функции, которая никогда ранее не возвращала `nil`, начать его возвращать.

Разумеется, никто не мешает вызывать функцию

```go
func foo(ctx context.Context, user User, order Order) (Receipt, error) {
    // ...
}
```

как `foo(context.TODO(), User{}, Order{})`, однако это хотя бы сообщает, что в функцию должны быть переданы non-nil значения, а валидация отдельных полей - уже ответственность самой функции. Однако уже можно быть уверенным, что при доступе к `user` и `order` паник не будет, и явные проверки типа `if order == nil { return }` уже не нужны.

### Возвращаемое значение

При возврате функцией значения вместо указателя, написание `return Receipt{}, err` занимает лишь немногим больше времени, чем `return nil, err`, однако у вызывающей стороны не будет необходимости проверять значение на `nil` (кто из нас не делал `return nil, nil` хотя бы раз в жизни?), не будет необходимости разыменовывать указатель при передаче куда-либо. Преимущества тут уже не настолько видны, т. к. принято возвращать либо значение и ошибку `nil`, либо значение `nil`/пустое значение и ошибку.

## Производительность

Бывают ситуации, когда функции принимают и возвращают указатели на большие по мнению разработчиков структуры, поэтому и передача их осуществляется по указателю с целью сэкономить время на копировании структуры. При этом, как правило, никто не пишет бенчмарки и не смотрит аннотации компилятора, потому что ситуация кажется очевидной - зачем передавать структуру размером в 1 КБ, когда можно передать указатель в 128 раз меньше?

Однако на практике данное решение может привести к менее производительному коду. Например, если взять маппинг двух структур:

```go
type User struct {
	ID                              int64
	CreatedAt, UpdatedAt, DeletedAt time.Time

	FirstName, SecondName, Patronymic string
	Birthday                          time.Time
	Nationality                       string
	UserType                          int

	Balance     *big.Rat
	BonusPoints *big.Rat
}

type UserDTO struct {
	ID                              int64
	CreatedAt, UpdatedAt, DeletedAt time.Time

	FirstName, SecondName, Patronymic string
	Birthday                          time.Time
	Nationality                       string
	UserType                          int

	Balance     *big.Rat
	BonusPoints *big.Rat

	FullName               string
	BalanceWithBonusPoints *big.Rat
}

func UserToDTO(u User) UserDTO {
	return UserDTO{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,

		FirstName:   u.FirstName,
		SecondName:  u.SecondName,
		Patronymic:  u.Patronymic,
		Birthday:    u.Birthday,
		Nationality: u.Nationality,
		UserType:    u.UserType,

		Balance:     u.Balance,
		BonusPoints: u.BonusPoints,

		FullName:               u.FirstName + " " + u.SecondName + " " + u.Patronymic,
		BalanceWithBonusPoints: new(big.Rat).Add(u.Balance, u.BonusPoints),
	}
}

func DTOToUser(d UserDTO) User {
	return User{
		ID:        d.ID,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
		DeletedAt: d.DeletedAt,

		FirstName:   d.FirstName,
		SecondName:  d.SecondName,
		Patronymic:  d.Patronymic,
		Birthday:    d.Birthday,
		Nationality: d.Nationality,
		UserType:    d.UserType,

		Balance:     d.Balance,
		BonusPoints: d.BonusPoints,
	}
}

func UserPtrToDTO(u *User) *UserDTO {
	return &UserDTO{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,

		FirstName:   u.FirstName,
		SecondName:  u.SecondName,
		Patronymic:  u.Patronymic,
		Birthday:    u.Birthday,
		Nationality: u.Nationality,
		UserType:    u.UserType,

		Balance:     u.Balance,
		BonusPoints: u.BonusPoints,

		FullName:               u.FirstName + " " + u.SecondName + " " + u.Patronymic,
		BalanceWithBonusPoints: new(big.Rat).Add(u.Balance, u.BonusPoints),
	}
}

func DTOPtrToUser(d *UserDTO) *User {
	return &User{
		ID:        d.ID,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
		DeletedAt: d.DeletedAt,

		FirstName:   d.FirstName,
		SecondName:  d.SecondName,
		Patronymic:  d.Patronymic,
		Birthday:    d.Birthday,
		Nationality: d.Nationality,
		UserType:    d.UserType,

		Balance:     d.Balance,
		BonusPoints: d.BonusPoints,
	}
}
```

и такой бенчмарк:

```go
func BenchmarkMapValues(b *testing.B) {
	var (
		user = createUser()
		res  pointers.User
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dto := pointers.UserToDTO(user)
		res = pointers.DTOToUser(dto)
	}

	_ = res
}

func BenchmarkMapPointers(b *testing.B) {
	var (
		user = createUser()
		res  *pointers.User
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dto := pointers.UserPtrToDTO(&user)
		res = pointers.DTOPtrToUser(dto)
	}

	_ = res
}

func createUser() pointers.User {
	return pointers.User{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: time.Now(),

		FirstName:   "John",
		SecondName:  "Doe",
		Patronymic:  "Smith",
		Birthday:    time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		Nationality: "Russian",
		UserType:    1,

		Balance:     big.NewRat(1000, 1),
		BonusPoints: big.NewRat(100, 1),
	}
}
```

и запустить его вот так:

```shell
go test -benchmem -bench . ./...
```

то можно получить вот такой результат:

```shell
goos: windows
goarch: amd64
pkg: articles/src/pointers
cpu: AMD Ryzen 7 5800H with Radeon Graphics
BenchmarkMapValues-16            4174956               289.7 ns/op           288 B/op          8 allocs/op
BenchmarkMapPointers-16          3207511               385.7 ns/op           704 B/op         10 allocs/op
```

Структура, разумеется, небольшая, но что более важно - при каждом вызове функции, возвращающей указатель, происходит дополнительное выделение памяти в куче, чего нет в функции, возвращающей значение.

Уберем из функций операции, приводящие к дополнительному выделению памяти в куче:

<details>
<summary>Код</summary>

```go
type User_NoHeap struct {
	ID                              int64
	CreatedAt, UpdatedAt, DeletedAt time.Time

	FirstName, SecondName, Patronymic string
	Birthday                          time.Time
	Nationality                       string
	UserType                          int

	Balance     *big.Rat
	BonusPoints *big.Rat
}

type UserDTO_NoHeap struct {
	ID                              int64
	CreatedAt, UpdatedAt, DeletedAt time.Time

	FirstName, SecondName, Patronymic string
	Birthday                          time.Time
	Nationality                       string
	UserType                          int

	Balance     *big.Rat
	BonusPoints *big.Rat
}

func UserToDTO_NoHeap(u User_NoHeap) UserDTO_NoHeap {
	return UserDTO_NoHeap{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,

		FirstName:   u.FirstName,
		SecondName:  u.SecondName,
		Patronymic:  u.Patronymic,
		Birthday:    u.Birthday,
		Nationality: u.Nationality,
		UserType:    u.UserType,

		Balance:     u.Balance,
		BonusPoints: u.BonusPoints,
	}
}

func DTOToUser_NoHeap(d UserDTO_NoHeap) User_NoHeap {
	return User_NoHeap{
		ID:        d.ID,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
		DeletedAt: d.DeletedAt,

		FirstName:   d.FirstName,
		SecondName:  d.SecondName,
		Patronymic:  d.Patronymic,
		Birthday:    d.Birthday,
		Nationality: d.Nationality,
		UserType:    d.UserType,

		Balance:     d.Balance,
		BonusPoints: d.BonusPoints,
	}
}

func UserPtrToDTO_NoHeap(u *User_NoHeap) *UserDTO_NoHeap {
	return &UserDTO_NoHeap{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,

		FirstName:   u.FirstName,
		SecondName:  u.SecondName,
		Patronymic:  u.Patronymic,
		Birthday:    u.Birthday,
		Nationality: u.Nationality,
		UserType:    u.UserType,

		Balance:     u.Balance,
		BonusPoints: u.BonusPoints,
	}
}

func DTOPtrToUser_NoHeap(d *UserDTO_NoHeap) *User_NoHeap {
	return &User_NoHeap{
		ID:        d.ID,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
		DeletedAt: d.DeletedAt,

		FirstName:   d.FirstName,
		SecondName:  d.SecondName,
		Patronymic:  d.Patronymic,
		Birthday:    d.Birthday,
		Nationality: d.Nationality,
		UserType:    d.UserType,

		Balance:     d.Balance,
		BonusPoints: d.BonusPoints,
	}
}
```

</details>

Также возьмем структуры размером 2 КБ, 8 КБ и 1 МБ (функции также без дополнительного выделения памяти в куче):

```go
type User2KB struct {
	Data [2048]byte
}

type UserDTO2KB struct {
	Data [2048]byte
}

func UserToDTO2KB(u User2KB) UserDTO2KB {
	return UserDTO2KB{
		Data: u.Data,
	}
}

func DTOToUser2KB(d UserDTO2KB) User2KB {
	return User2KB{
		Data: d.Data,
	}
}

func UserPtrToDTO2KB(u *User2KB) *UserDTO2KB {
	return &UserDTO2KB{
		Data: u.Data,
	}
}

func DTOPtrToUser2KB(d *UserDTO2KB) *User2KB {
	return &User2KB{
		Data: d.Data,
	}
}
```

Результат при этом получается вот таким:

```shell
goos: windows
goarch: amd64
pkg: articles/src/pointers
cpu: AMD Ryzen 7 5800H with Radeon Graphics
BenchmarkMapValues1MB-16                    9324            123896 ns/op               0 B/op          0 allocs/op
BenchmarkMapPointers1MB-16                  5163            218242 ns/op         2097156 B/op          2 allocs/op
BenchmarkMapValues2KB-16                 6429981               186.2 ns/op             0 B/op          0 allocs/op
BenchmarkMapPointers2KB-16               2866360               416.5 ns/op          2048 B/op          1 allocs/op
BenchmarkMapValues8KB-16                 2202045               544.4 ns/op             0 B/op          0 allocs/op
BenchmarkMapPointers8KB-16                773160              1548 ns/op            8192 B/op          1 allocs/op
BenchmarkMapValuesNoHeap-16             41105602                28.50 ns/op            0 B/op          0 allocs/op
BenchmarkMapPointersNoHeap-16           18602142                58.65 ns/op          192 B/op          1 allocs/op
BenchmarkMapValues-16                    3729990               330.7 ns/op           288 B/op          8 allocs/op
BenchmarkMapPointers-16                  2999234               396.9 ns/op           704 B/op         10 allocs/op
```

Также я взял большую структуру из кода на работе (которую не покажу), которая по размеру немного меньше 1 КБ. Результаты замены единственной строки (`func convert(x *X) *Y` -> `func convert(x X) Y`) таковы (третий бенчмарк - это передача `*&X{}`, а в остальном он аналогичен первому):

```shell
goos: windows
goarch: amd64
pkg: supercompany/code
cpu: AMD Ryzen 7 5800H with Radeon Graphics
BenchmarkValue-16                                2789080               411.7 ns/op           432 B/op          7 allocs/op
BenchmarkPointer-16                              1767438               673.8 ns/op          1136 B/op          8 allocs/op
BenchmarkDereferencePointerToValue-16            2450964               422.3 ns/op           432 B/op          7 allocs/op
```

### Fizzbuzz

Как можно написать серьезную статью без как минимум одного сложного и научного алгоритма? Реализуем fizzbuzz!

Сделаем две реализации: одна будет по возможности использовать передачу структур по значению (копирование), а вторая - передачу по указателю, а затем сравним их производительность.

<details>
<summary>Реализация</summary>

```go
type ControllerReq struct {
	From, To int
}

type ControllerResp struct {
	Values map[int]string
}

type logicReq struct {
	value int
}

type logicResp struct {
	value string
}

func ValueController(ctx context.Context, req *ControllerReq) (*ControllerResp, error) {
	res := make(map[int]string, req.To-req.From)
	for i := req.From; i < req.To; i++ {
		x, err := valueLogic(ctx, logicReq{i})
		if err != nil {
			return nil, err
		}

		res[i] = x.value
	}

	return &ControllerResp{res}, nil
}

func valueLogic(ctx context.Context, req logicReq) (logicResp, error) {
	var (
		divisibleBy3 = req.value%3 == 0
		divisibleBy5 = req.value%5 == 0
	)
	switch {
	case divisibleBy3 && divisibleBy5:
		return logicResp{"fizzbuzz"}, nil
	case divisibleBy3:
		return logicResp{"fizz"}, nil
	case divisibleBy5:
		return logicResp{"buzz"}, nil
	default:
		return logicResp{strconv.FormatInt(int64(req.value), 10)}, nil
	}
}

func PtrController(ctx context.Context, req *ControllerReq) (*ControllerResp, error) {
	res := make(map[int]string, req.To-req.From)
	for i := req.From; i < req.To; i++ {
		x, err := ptrLogic(ctx, &logicReq{i})
		if err != nil {
			return nil, err
		}

		res[i] = x.value
	}

	return &ControllerResp{res}, nil
}

func ptrLogic(ctx context.Context, req *logicReq) (*logicResp, error) {
	var (
		divisibleBy3 = req.value%3 == 0
		divisibleBy5 = req.value%5 == 0
	)
	switch {
	case divisibleBy3 && divisibleBy5:
		return &logicResp{"fizzbuzz"}, nil
	case divisibleBy3:
		return &logicResp{"fizz"}, nil
	case divisibleBy5:
		return &logicResp{"buzz"}, nil
	default:
		return &logicResp{strconv.FormatInt(int64(req.value), 10)}, nil
	}
}
```

</details>

<details>
<summary>Бенчмарк</summary>

```go
func BenchmarkValue(b *testing.B) {
	var (
		req  = &fizzbuzz.ControllerReq{From: 1, To: 100}
		resp *fizzbuzz.ControllerResp
		err  error
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err = fizzbuzz.ValueController(context.TODO(), req)
	}

	_, _ = resp, err
}

func BenchmarkPointer(b *testing.B) {
	var (
		req  = &fizzbuzz.ControllerReq{From: 1, To: 100}
		resp *fizzbuzz.ControllerResp
		err  error
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err = fizzbuzz.PtrController(context.TODO(), req)
	}

	_, _ = resp, err
}
```

</details>

Результат:

```shell
goos: windows
goarch: amd64
pkg: articles/src/fizzbuzz
cpu: AMD Ryzen 7 5800H with Radeon Graphics
BenchmarkValue-16         323662              3313 ns/op            4227 B/op          4 allocs/op
BenchmarkPointer-16       204608              5862 ns/op            5811 B/op        103 allocs/op
```

Как можно видеть, реализация с передачей по значению выделяет намного меньше памяти, делает меньше аллокаций и работает быстрее. Разумеется, это выдуманный пример, но реализовывать огромное приложение для демонстрации разницы производительности было бы непрактично.

## Расширяемость

Под расширяемостью я понимаю способность кода не требовать модификаций при изменении требований или связанного с ним кода.

Например, если мы добавим в структуру мьютекс, то передавать его копии будет уже невозможно, и в таком случае старый код может потребоваться переписать. С другой стороны, большинство обрабатываемых приложением данных не требуют конкурентного доступа, а оттого, как правило, не будут содержать поля, не позволяющие копирование.

## Консистентность

Практически любой репозиторий обязательно будет иметь параметры или возвращаемые значения, являющимися как указателями, так и значениями. Поэтому говорить о желании сделать код более единообразным не получится - целые числа и строки всё равно будут передаваться по значению, а структуры с мьютексами - по указателю.

Наоборот, передача параметров и возврат значений двумя способами могут повысить читаемость кода, т. к. будут обусловлены особенностями передаваемой структуры и её использованием (наличием мьютекса, необходимостью передавать по указателю для изменения или сохранения указателя на значение), а не принятым в команде соглашением всё передавать по указателю. Передача по указателю будет выделяться и иметь необходимость, вместо того чтобы быть выбором по умолчанию.

## Потенциальные ошибки

Повсеместное использование указателей в качестве параметров приводит к выбору: необходимо либо проверять каждый параметр на равенство `nil`, либо допускать, что произойдёт паника при попытке разыменования указателя `nil`.

Передача по значению может привести к случайному копированию и изменению значений полей у копии, а не у оригинального значения, но такие вещи легко обнаруживаются линтерами, на ревью и здравым смыслом. Не утверждаю, что такие ошибки исключены полностью, но _для меня_ они всегда оказывались более редкими, менее фатальными и легко обнаруживались второй парой глаз.
