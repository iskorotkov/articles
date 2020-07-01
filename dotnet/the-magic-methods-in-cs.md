Магические сигнатуры методов в C\#

[CEZARY PIĄTEK](https://cezarypiatek.github.io/post/methods-with-special-signature/)

Представляю вашему вниманию перевод статьи [The Magical Methods in C#](https://cezarypiatek.github.io/post/methods-with-special-signature/) автора [CEZARY PIĄTEK](https://cezarypiatek.github.io/).

Есть определенный набор сигнатур методов в C#, имеющих поддержку на уровне языка. Методы с такими сигнатурами позволяют использовать специальный синтаксис со всеми его преимуществами. Например, с их помощью можно упростить наш код или создать DSL для того, чтобы выразить решение проблемы более красивым образом. Я встречаюсь с такими методами повсеместно, так что я решил написать пост и обобщить все мои находки по этой теме, а именно:

- Синтаксис инициализации коллекций
- Синтаксис инициализации словарей
- Деконструкторы
- Пользовательские awaitable типы
- Паттерн query expression

## Синтаксис инициализации коллекций

[Синтаксис инициализации коллекции](https://docs.microsoft.com/en-us/dotnet/csharp/programming-guide/classes-and-structs/object-and-collection-initializers#collection-initializers) довольно старая фича, т. к. она существует с C# 3.0 (выпущен в конце 2007 года). Напомню, синтаксис инициализации коллекции позволяет создать список с элементами в одном блоке:

```cs
var list = new List<int> { 1, 2, 3 };
```

Этот код эквивалентен приведенному ниже:

```cs
var temp = new List<int>();
temp.Add(1);
temp.Add(2);
temp.Add(3);
var list = temp;
```

Возможность использования синтаксиса инициализации коллекции не ограничивается только классами из BCL. Он может быть использован с любым типом, удовлетворяющим следующим условиям:

- тип имплементирует интерфейс `IEnumerable`
- тип имеет метод с сигнатурой `void Add(T item)`

```cs
public class CustomList<T>: IEnumerable
{
    public IEnumerator GetEnumerator() => throw new NotImplementedException();
    public void Add(T item) => throw new NotImplementedException();
}
```

Мы можем добавить поддержку синтаксиса инициализации коллекции, определив `Add` как метод расширения:

```cs
public static class ExistingTypeExtensions
{
    public static void Add<T>(this ExistingType @this, T item) => throw new NotImplementedException();
}
```

Этот синтаксис также можно использовать для вставки элементов в поле-коллекцию без публичного сеттера:

```cs
class CustomType
{
    public List<string> CollectionField { get; private set; }  = new List<string>();
}

class Program
{
    static void Main(string[] args)
    {
        var obj = new CustomType
        {
            CollectionField =
            {
                "item1",
                "item2"
            }
        };
    }
}
```

Синтаксис инициализации коллекции полезен при инициализации коллекции известным числом элементов. Но что если мы хотим создать коллекцию с переменным числом элементов? Для этого есть менее известный синтаксис:

```cs
var obj = new CustomType
{
    CollectionField =
    {
        { existingItems }
    }
};
```

Такое возможно для типов, удовлетворяющих следующим условиям:

- тип имплементирует интерфейс `IEnumerable`
- тип имеет метод с сигнатурой `void Add(IEnumerable<T> items)`

```cs
public class CustomList<T>: IEnumerable
{
    public IEnumerator GetEnumerator() => throw new NotImplementedException();
    public void Add(IEnumerable<T> items) => throw new NotImplementedException();
}
```

К сожалению, массивы и коллекции из BCL не реализуют метод `void Add(IEnumerable<T> items)`, но мы можем изменить это, определив метод расширения для существующих типов коллекций:

```cs
public static class ListExtensions
{
    public static void Add<T>(this List<T> @this, IEnumerable<T> items) => @this.AddRange(items);
}
```

Благодаря этому мы можем написать следующее:

```cs
var obj = new CustomType
{
    CollectionField =
    {
        { existingItems.Where(x => /*Filter items*/).Select(x => /*Map items*/) }
    }
};
```

Или даже собрать коллекцию из смеси индивидуальных элементов и результатов нескольких перечислений (IEnumerable):

```cs
var obj = new CustomType
{
    CollectionField =
    {
        individualElement1,
        individualElement2,
        { list1.Where(x => /*Filter items*/).Select(x => /*Map items*/) },
        { list2.Where(x => /*Filter items*/).Select(x => /*Map items*/) },
    }
};
```

Без подобного синтаксиса очень сложно получить подобный результат в блоке инициализации.

Я узнал об этой фиче совершенно случайно, когда работал с маппингами для типов с полями-коллекциями, сгенерированными из контрактов `protobuf`. Для тех, кто не знаком с `protobuf`: если вы используете [grpctools](https://developers.google.com/protocol-buffers/docs/reference/csharp-generated) для генерации типов .NET из файлов `proto`, все поля-коллекции генерируются подобным образом:

```cs
[DebuggerNonUserCode]
public RepeatableField<ItemType> SomeCollectionField
{
    get
    {
        return this.someCollectionField_;
    }
}
```

Как можно заметить, поля-коллекции не имеют сеттер, но [`RepeatableField`](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/collections/repeated-field-t-) реализует метод [`void Add(IEnumerable items)`](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/collections/repeated-field-t-#class_google_1_1_protobuf_1_1_collections_1_1_repeated_field_3_01_t_01_4_1a6d0e9efbac818182068afae48b8d4599), так что мы по-прежнему можем инициализировать их в блоке инициализации:

```cs
/// <summary>
/// Adds all of the specified values into this collection. This method is present to
/// allow repeated fields to be constructed from queries within collection initializers.
/// Within non-collection-initializer code, consider using the equivalent <see cref="AddRange"/>
/// method instead for clarity.
/// </summary>
/// <param name="values">The values to add to this collection.</param>
public void Add(IEnumerable<T> values)
{
    AddRange(values);
}
```

## Синтаксис инициализации словарей

Одна из крутых фич C# 6.0 - [инициализация словаря по индексу](https://docs.microsoft.com/en-us/dotnet/csharp/whats-new/csharp-6#initialize-associative-collections-using-indexers), которая упростила синтаксис инициализации словарей. Благодаря ей мы можем писать более читаемый код:

```cs
var errorCodes = new Dictionary<int, string>
{
    [404] = "Page not Found",
    [302] = "Page moved, but left a forwarding address.",
    [500] = "The web server can't come out to play today."
};
```

Этот код эквивалентен следующему:

```cs
var errorCodes = new Dictionary<int, string>();
errorCodes[404] = "Page not Found";
errorCodes[302] = "Page moved, but left a forwarding address.";
errorCodes[500] = "The web server can't come out to play today.";
```

Это немного, но это определенно упрощает написание и чтение кода.

Лучшее в инициализации по индексу - это то, что она не ограничивается классом `Dictionary<T>` и может быть использована с любым другим типом, определившим индексатор:

```cs
class HttpHeaders
{
    public string this[string key]
    {
        get => throw new NotImplementedException();
        set => throw new NotImplementedException();
    }
}

class Program
{
    static void Main(string[] args)
    {
        var headers = new HttpHeaders
        {
            ["access-control-allow-origin"] = "*",
            ["cache-control"] = "max-age=315360000, public, immutable"
        };
    }
}
```

## Деконструкторы

В C# 7.0 помимо кортежей был добавлен механизм деконструкторов. Они позволяют декомпозировать кортеж в набор отдельных переменных:

```cs
var point = (5, 7);
// decomposing tuple into separated variables
var (x, y) = point;
```

Что эквивалентно следующему:

```cs
ValueTuple<int, int> point = new ValueTuple<int, int>(1, 4);
int x = point.Item1;
int y = point.Item2;
```

Этот синтаксис позволяет обменять значения двух переменных без явного объявления третьей:

```cs
int x = 5, y = 7;
//switch
(x, y) = (y,x);
```

Или использовать более краткий метод инициализации членов класса:

```cs
class Point
{
    public int X { get; }
    public int Y { get; }

    public Point(int x, int y)  => (X, Y) = (x, y);
}
```

Деконструкторы могут быть использованы не только с кортежами, но и с другими типами. Для использования деконструкции типа этот тип должен реализовывать метод, подчиняющийся следующим правилам:

- метод называется `Deconstruct`
- метод возвращает `void`
- все параметры метода имеют модификатор `out`

Для нашего типа `Point` мы можем объявить деконструктор следующим образом:

```cs
class Point
{
    public int X { get; }
    public int Y { get; }

    public Point(int x, int y) => (X, Y) = (x, y);

    public void Deconstruct(out int x, out int y) => (x, y) = (X, Y);
}
```

Пример использования приведен ниже:

```cs
var point = new Point(2, 4);
var (x, y) = point;
```

"Под капотом" он превращается в следующее:

```cs
int x;
int y;
new Point(2, 4).Deconstruct(out x, out y);
```

Деконструкторы могут быть добавлены к типам с помощью методов расширения:

```cs
public static class PointExtensions
{
     public static void Deconstruct(this Point @this, out int x, out int y) => (x, y) = (@this.X, @this.Y);
}
```

Один из самых полезных примеров применения деконструкторов - это деконструкция `KeyValuePair<TKey, TValue>`, которая позволяет с легкостью получить доступ к ключу и значению во время итерирования по словарю:

```cs
foreach (var (key, value) in new Dictionary<int, string> { [1] = "val1", [2] = "val2" })
{
    //TODO: Do something
}
```

`KeyValuePair<TKey, TValue>.Deconstruct(TKey, TValue)` доступно только с `netstandard2.1`. Для предыдущих версий `netstandard` нам нужно использовать ранее приведенный метод расширения.

## Пользовательские awaitable типы

C# 5.0 (выпущен вместе с Visual Studio 2012) ввел механизм `async/await`, который стал переворотом в области асинхронного программирования. Прежде вызов асинхронного метода представлял собой запутанный код, особенно когда таких вызовов было несколько:

```cs
void DoSomething()
{
    DoSomethingAsync().ContinueWith((task1) => {
        if (task1.IsCompletedSuccessfully)
        {
            DoSomethingElse1Async(task1.Result).ContinueWith((task2) => {
                if (task2.IsCompletedSuccessfully)
                {
                    DoSomethingElse2Async(task2.Result).ContinueWith((task3) => {
                        //TODO: Do something
                    });
                }
            });
        }
    });
}

private Task<int> DoSomethingAsync() => throw new NotImplementedException();
private Task<int> DoSomethingElse1Async(int i) => throw new NotImplementedException();
private Task<int> DoSomethingElse2Async(int i) => throw new NotImplementedException();
```

Это может быть переписано намного красивее с использованием синтаксиса `async/await`:

```cs
async Task DoSomething()
{
    var res1 = await DoSomethingAsync();
    var res2 = await DoSomethingElse1Async(res1);
    await DoSomethingElse2Async(res2);
}
```

Это может прозвучать удивительно, но ключевое слово `await` не зарезервировано только под использование с типом `Task`. Оно может быть использовано с любым типом, который имеет метод `GetAwaiter`, возвращающий удовлетворяющий следующим требованиям тип:

- тип имплементирует интерфейс `System.Runtime.CompilerServices.INotifyCompletion` и реализует метод `void OnCompleted(Action continuation)`
- тип имеет свойство `IsCompleted` логического типа
- тип имеет метод `GetResult` без параметров

Для добавления поддержки ключевого слова `await` к пользовательскому типу мы должны определить метод `GetAwaiter`, возвращающий `TaskAwaiter<TResult>` или пользовательский тип, удовлетворяющий приведенным выше условиям:

```cs
class CustomAwaitable
{
    public CustomAwaiter GetAwaiter() => throw new NotImplementedException();
}

class CustomAwaiter: INotifyCompletion
{
    public void OnCompleted(Action continuation) => throw new NotImplementedException();

    public bool IsCompleted => throw new NotImplementedException();

    public void GetResult() => throw new NotImplementedException();
}
```

Вы можете спросить: "Каков возможный сценарий использования синтаксиса `await` с пользовательским awaitable типом?". Если это так, то я рекомендую вам прочитать статью [Stephen Toub](https://devblogs.microsoft.com/pfxteam/author/toub/) под названием "[await anything](https://devblogs.microsoft.com/pfxteam/await-anything/)", которая показывает множество интересных примеров.

## Паттерн query expression

Лучшее нововведение C# 3.0 - Language-Integrated Query, также известное как LINQ, предназначенное для манипулирования коллекциями с SQL-подобным синтаксисом. LINQ имеет две вариации: SQL-подобный синтаксис и синтаксис методов расширения. Я предпочитаю второй вариант, т. к. по моему мнению он более читаем, а также потому что я привык к нему. Интересный факт о LINQ заключается в том, что SQL-подобный синтаксис во время компиляции транслируется в синтаксис методов расширения, т. к. это фича C#, а не CLR. LINQ был разработан в первую очередь для работы с типами `IEnumerable`, `IEnumerable<T>` и `IQuerable<T>`, но он не ограничен только ими, и мы можем использовать его с любым типом, удовлетворяющим требованиям [паттерна query expression](https://github.com/dotnet/csharplang/blob/master/spec/expressions.md#the-query-expression-pattern). Полный набор сигнатур методов, используемых LINQ, таков:

```cs
class C
{
    public C<T> Cast<T>();
}

class C<T> : C
{
    public C<T> Where(Func<T,bool> predicate);

    public C<U> Select<U>(Func<T,U> selector);

    public C<V> SelectMany<U,V>(Func<T,C<U>> selector, Func<T,U,V> resultSelector);

    public C<V> Join<U,K,V>(C<U> inner, Func<T,K> outerKeySelector, Func<U,K> innerKeySelector, Func<T,U,V> resultSelector);

    public C<V> GroupJoin<U,K,V>(C<U> inner, Func<T,K> outerKeySelector, Func<U,K> innerKeySelector, Func<T,C<U>,V> resultSelector);

    public O<T> OrderBy<K>(Func<T,K> keySelector);

    public O<T> OrderByDescending<K>(Func<T,K> keySelector);

    public C<G<K,T>> GroupBy<K>(Func<T,K> keySelector);

    public C<G<K,E>> GroupBy<K,E>(Func<T,K> keySelector, Func<T,E> elementSelector);
}

class O<T> : C<T>
{
    public O<T> ThenBy<K>(Func<T,K> keySelector);

    public O<T> ThenByDescending<K>(Func<T,K> keySelector);
}

class G<K,T> : C<T>
{
    public K Key { get; }
}
```

Разумеется, мы не обязаны реализовывать все эти методы для того, чтобы использовать синтаксис `LINQ` с нашим пользовательским типом. Список обязательных операторов и методов `LINQ` для них можно посмотреть [здесь](https://docs.microsoft.com/en-us/dotnet/csharp/programming-guide/concepts/linq/query-expression-syntax-for-standard-query-operators). Действительно хорошее объяснение того, как это сделать, можно найти в статье [Understand monads with LINQ](https://codewithstyle.info/understand-monads-linq/) автора [Miłosz Piechocki](http://miloszpiechocki.com/).

## Подведение итогов

Цель этой статьи заключается вовсе не в том, чтобы убедить вас злоупотреблять этими синтаксическими трюками, а в том, чтобы сделать их более понятными. С другой стороны, их нельзя всегда избегать. Они были разработаны для того, чтобы их использовать, и иногда они могут сделать ваш код лучше. Если вы боитесь, что получившийся код будет непонятен вашим коллегам, вам нужно найти способ поделиться знаниями с ними (или хотя бы ссылкой на эту статью). Я не уверен, что это полный набор таких "магических методов", так что если вы знаете еще какие-то - пожалуйста, поделитесь в комментариях.
