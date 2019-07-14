Source: https://nora.codes/post/what-is-rusts-unsafe/

# Что значит unsafe в Rust?

Мне доводилось видеть много недопониманий относительно того, что значит ключевое слово unsafe для полезности и правильности языка Rust и его продвижения как "безопасного языка системного программирования". Правда намного сложнее, чем можно описать в  коротком твите, к сожалению. Вот как я ее вижу.

В целом, **ключевое слово unsafe не выключает систему типов, которая поддерживает код на Rust корректным**. Она только дает возможность использовать некоторые "суперспособности", такие как разыменование указателей. unsafe используется для реализации безопасных абстракций на основе фундаментально небезопасного мира, чтобы большая часть кода на Rust могла использовать эти абстракции и избегать небезопасного доступа к памяти.

# Гарантия безопасности

Rust гарантирует безопасность как один из своих главных принципов. Можно сказать, что это *смысл существования* языка. Он, однако, не предоставляет безопасность в традиционном смысле, во время выполнения программы и с использованием сборщика мусора. Вместо этого Rust использует очень продвинутую систему типов, чтобы следить за тем, когда и к каким значениям можно получать доступ. Затем компилятор статически анализирует каждую программу на Rust, чтобы убедиться в том, что она постоянно находится в корректном состоянии.

## Безопасность в Python

Давайте возьмем для примера Python. Чистый код на Python не может повредить память. Доступ к элементам списков имеет проверки на выход за границы; ссылки, возвращаемые функциями подсчитываются во избежание появления висячих ссылок; нет никакого способа производить произвольные арифметические операции с указателями.

У этого есть два следствия. Во-первых, много типов должны быть "специальными". Например, невозможно реализовать эффективный список или словарь на чистом Python. Вместо этого интерпретатор CPython имеет их внутреннюю реализацию. Во-вторых, доступ к внешним функциям (функциям, реализованным не на Python), называемый интерфейсом внешней функции, требует использования специального модуля ctypes и нарушает гарантии безопасности языка.

В каком-то смысле это значит, что все, что написано на Python, не гарантирует безопасного доступа к памяти.

## Безопасность в Rust

Rust тоже предоставляет безопасность, но вместо реализации небезопасных структур на C, он предоставляет уловку: ключевое слово unsafe. Это означает, что фундаментальные структуры данных в Rust, такие как Vec, VecDeque, BTreeMap и String, реализованы на языке Rust.

Вы спросите: "Но если Rust предоставляет уловку против его гарантий безопасности кода, и стандартная библиотека реализована с использованием этой уловки, разве не все в Rust будет считаться небезопасным?"

Одним словом, уважаемый читатель, - **да**, именно так, как это было в Python. Давайте разберем это подробнее.

# Что запрещено в безопасном Rust?

Безопасность в Rust хорошо определена: мы много о ней думаем. Вкратце, безопасные программы на Rust не могут:

- **Разыменовывать указатель, указывающий не на тот тип, о котором знает компилятор**. Это значит, что не существует никаких указателей на null (потому что они никуда не указывают), никаких ошибок выхода за границы и/или ошибок сегментации (segmentation faults), никаких переполнений буфера. Но также это значит и то, что нет никаких использований после освобождения памяти или повторного освобождения памяти (потому что освобождение памяти считается разыменованием указателя) и никакого [каламбура типизации](https://en.wikipedia.org/wiki/Type_punning).

- **Иметь несколько изменяемых ссылок на объект или одновременно изменяемые и неизменяемые ссылки на объект**. То есть если у вас есть изменяемая ссылка на объект, у вас может быть только она, а если у вас есть неизменяемая ссылка на объект, он не изменится, пока она у вас сохраняется. Это означает невозможность вызвать гонку данных в безопасном Rust, что является гарантией, которую большинство других безопасных языков предоставить не могут.

Rust кодирует эту информацию в системе типов или используя **алгебраические типы данных**, такие как Option<T> для обозначения существования/отсутствия значения и Result<T, E> для обозначения ошибки/успеха, или **ссылки и их время жизни**, например, &T vs &mut T для обозначения общей (неизменяемой) ссылки и эксклюзивной (изменяемой) ссылки и &'a T vs &'b T для различения ссылок, которые являются корректными в различных контекстах (такое, как правило, опускается, так как компилятор достаточно умный, чтобы сам понять это).

## Примеры

Например, следующий код не будет компилироваться, так как он содержит висячую ссылку. Конкретнее, *my_struct does not live enough*. Иными словами, функция вернет ссылку на что-то, что уже не существует, и поэтому компилятор не сможет (и, на самом деле, даже не знает, как) скомпилировать это.

```rust
fn dangling_reference(v: &u64) -> &MyStruct {
    // Создаем новое значение типа MyStruct со значением поля, равным v, единственным параметром функции.
    let my_struct = MyStruct { value: v };
    // Возвращаем ссылку на локальную переменную my_struct.
    return &my_struct;
    // Память из-под my_struct освобождается (снимается со стека).
}
```

Этот код делает то же самое, но он пытается обойти эту проблему, размещая значение в куче (Box - это имя базового умного указателя в Rust).

```rust
fn dangling_heap_reference(v: &u64) -> &Box<MyStruct> {
    let my_struct = MyStruct { value: v };
    // Помещаем структуру в Box с выделением места в куче и перемещением ее туда.
    let my_box = Box::new(my_struct);
    // Возвращаем ссылку на локальную переменную my_box.
    return &my_box;
    // my_box снимается со стека. Эта переменная "владеет" my_struct и поэтому ответственна за освобождение памяти из-под нее,
    // так что память из-под MyStruct тоже освобождается.
}
```

Правильный код возвращает сам Box<MyStruct> вместо ссылки на него. Это кодирует перемещение владения - ответственности за освобождение памяти - в сигнатуре функции. При взгляде на сигнатуру становится понятно, что вызывающий код ответственен за то, что произойдет с Box<MyStruct>, и, действительно, компилятор обрабатывает это автоматически.

```rust
fn no_dangling_reference(v: &u64) -> Box<MyStruct> {
    let my_struct = MyStruct { value: v };
    let my_box = Box::new(my_struct);
    // Возвращаем локальную переменную my_box по значению.
    return my_box;
    // Никакая память не освобождается. Вызывающий код теперь ответственен за управление памятью в куче,
    // выделенной в этой функции; она почти точно будет освобождена автоматически
    // когда Box<MyStruct> выйдет из области действия в вызывающем коде, за исключением случая возникновения двойной паники.
}
```

> Некоторые плохие вещи не запрещены в безопасном Rust. Например, разрешено с точки зрения компилятора:
>
> - вызвать deadlock в программе
> - совершить утечку произвольно большого объема памяти
> - не суметь закрыть хендлы файлов, соединения баз данных или крышки ракетных шахт
>
> Сила экосистемы Rust в том, что много проектов выбирают использование системы типов для обеспечения правильности кода по максимуму, но компилятор не требует такого принуждения, за исключением случаев обеспечения безопасного доступа к памяти.

# Что разрешено в небезопасном Rust?

Небезопасный код на Rust - это код на Rust с ключевым словом unsafe. unsafe может применяться к функции или блоку кода. Когда оно применяется к функции, это значит "эта функция требует, чтобы вызываемый код вручную обеспечивал инвариант, который обычно обеспечивается компилятором". Когда оно применяется к блоку кода, это значит "этот блок кода вручную обеспечивает инвариант, необходимый для предотвращения небезопасного доступа к памяти, и поэтому ему разрешено делать небезопасные вещи".

**Иными словами, unsafe у функции значит "тебе надо все проверить", а на блоке кода - "я уже все проверил".**

Как отмечено в книге [The Rust Programming Language](https://doc.rust-lang.org/book/ch19-01-unsafe-rust.html), код в блоке, отмеченном ключевым словом unsafe, может:

- **Разыменовывать указатель.** Это ключевая "суперспособность", которая позволяет реализовывать двусвязные списки, hashmap, и другие фундаментальные структуры данных.

- **Вызывать небезопасную функцию или метод.** Больше об этом ниже.

- **Получать доступ или модифицировать изменяемую статическую переменную.** Статические переменные, у которых область видимости не контролируется, не могут быть статически проверены, поэтому их использование небезопасно.

- **Имплементировать небезопасный типаж (trait).** Небезопасные типажи используются для того, чтобы помечать, гарантируют ли конкретные типы определенные инварианты. Например, Send и Sync определяют, может ли тип пересылаться между границами потоков или быть использован несколькими потоками одновременно.

Помните те примеры с висячими указателями, приведенные выше? Добавьте слово unsafe и компилятор станет ругаться вдвое больше, потому что ему не нравится использование unsafe там, где оно не требуется.

Вместо этого ключевое слово unsafe используется для реализации безопасных абстракций на основе произвольных операций над указателями. Например, тип Vec реализуется с использованием unsafe, но его безопасно использовать, так как он проверяет попытки получить доступ к элементам и не допускает переполнения. Хотя он и предоставляет операции наподобие set_len, которые *могут* вызвать небезопасность доступа к памяти, они помечены как unsafe.

Например, мы могли бы сделать то же самое, что и в примере no_dangling_reference, но с беспричинным использованием unsafe:

```rust
fn manual_heap_reference(v: u64) -> *mut MyStruct {
    let my_struct = MyStruct { value: v };
    let my_box = Box::new(my_struct);
    // Преобразовать Box в старый добрый указатель.
    let struct_pointer = Box::into_raw(my_box);
    return struct_pointer;
    // Ничего не разыменовывается; эта функция просто возвращает указатель.
    // MyStruct остается там же в куче.
}
```

Заметьте отсутствие слова unsafe. Создание указателей абсолютно безопасно. Как было написано, это риск возникновения утечки памяти, но ничего более, а утечки памяти безопасны. Вызов этой функции тоже безопасен. unsafe требуется только тогда, когда что-то пытается **разыменовать** указатель. Как дополнительный бонус разыменование позволит автоматически освободить выделенную память.

```rust
fn main() {
    let my_pointer = manual_heap_reference(1337);
    let my_boxed_struct = unsafe { Box::from_raw(my_pointer) };
    // Печатает "Value: 1337"
    println!("Value: {}", my_boxed_struct.value);
    // my_boxed_struct выходит из области видимости. Эта переменная теперь владеет памятью в куче, поэтому
    // она освобождает память из-под MyStruct
}
```

После оптимизации этот код эквивалентен простому возвращению Box. Box - безопасная абстракция на основе указателей, потому что она предотвращает распространение указателей повсюду. Например, следующая версия main **приведет** к двойному освобождению памяти (double-free).

```rust
fn main() {
    let my_pointer = manual_heap_reference(1337);
    let my_boxed_struct_1 = unsafe { Box::from_raw(my_pointer) };
    // DOUBLE FREE BUG!
    let my_boxed_struct_2 = unsafe { Box::from_raw(my_pointer) };
    // Печатает "Value: 1337" дважды.
    println!("Value: {}", my_boxed_struct_1.value);
    println!("Value: {}", my_boxed_struct_2.value);
    // my_boxed_struct_2 выходит из области видимости. Он владеет памятью в куче, поэтому
    // он освобождает память из-под MyStruct.
    // Затем my_boxed_struct_1 выходит из области видимости. Он также владеет памятью в куче,
    // поэтому он также освобождает память из-под MyStruct. Это double-free bug.
}
```

## Так что такое безопасная абстракция?

Безопасная абстракция - это такая абстракция, которая использует систему типов для предоставления API, который не может быть использован для нарушения гарантий безопасности, которые были упомянуты выше. Box безопаснее *mut T, так как он не может привести к двойному освобождению памяти, проиллюстрированному выше.

Другой пример - тип Rc в Rust. Это указатель с подсчетом ссылок - недопускающая изменение ссылка на данные, расположенные в куче. Так как она допускает множественный одновременный доступ к одной области памяти, она *должна* предотвращать изменение для того, чтобы считаться безопасной.

В дополнение к этому, она не потокобезопасна. Если вам нужна потокобезопасность, вам придется использовать тип Arc (Atomic Reference Counting), который имеет штраф к производительности из-за использования atomic значений для подсчета ссылок и предотвращения возможных гонок данных в многопоточных средах.

Компилятор не позволит вам использовать Rc там, где вы должны использовать Arc, потому что создатели типа Rc не отметили его как потокобезопасный. Если бы они сделали это, это было бы необоснованно: ложное обещание безопасности.

## Когда необходим небезопасный Rust?

Небезопасный Rust необходим всегда, когда необходимо произвести операцию, нарушающую одно из тех двух правил, описанных выше. Например, в двусвязном списке отсутствие изменяемых ссылок на одни и те же данные (у следующего элемента и предыдущего элемента) полностью лишает его пользы. С unsafe имплементатор двусвязного списка может написать код, используя указатели *mut Node<T> и затем инкапсулировать это в безопасную абстракцию.

Другой пример - работа со встраиваемыми системами. Часто микроконтроллеры используют набор регистров, чьи значения определяются физическим состоянием устройства. Мир не может остановиться, пока вы берете &mut u8 с такого регистра, поэтому для работы с крэйтами поддержки устройства необходим unsafe. Как правило, такие крэйты инкапсулируют состояние в прозрачных безопасных обертках, которые по возможности копируют данные, или же используют другие техники, обеспечивающие гарантии компилятора.

Иногда необходимо провести операцию, которая может привести к одновременному чтению и записи, или небезопасному доступу к памяти, и именно здесь нужен unsafe. Но до тех пор, пока есть возможность убедиться в поддержании безопасных инвариантов перед тем, как пользователь безопасного (то есть не отмеченного unsafe) кода коснется чего-то, все в порядке.

# На чьих плечах лежит эта ответственность?

Мы приходим к сделанному раньше утверждению - **да**, полезность кода на Rust основана на небезопасном коде. Несмотря на то, что это сделано несколько иначе, чем небезопасная имплементация основных структур данных в Python, реализация Vec, Hashmap и т. д. **должна** использовать манипуляции указателями в какой-либо степени.

Мы говорим, что Rust безопасен, с фундаментальным допущением, что небезопасный код, который мы используем через наши зависимости либо от стандартной библиотеки, либо от кода других библиотек, корректно написан и инкапсулирован. Фундаментальное преимущество Rust заключается в том, что небезопасный код загнан в небезопасные блоки, которые должны быть тщательно проверены их авторами.

В Python бремя проверки безопасности манипуляций с памятью лежит только на разработчиках интерпретаторов и пользователях интерфейсов внешних функций. В C это бремя лежит на каждом программисте.

В Rust оно лежит на пользователях ключевого слова unsafe. Это очевидно, так как внутри такого кода инварианты должны поддерживаться вручную, и поэтому необходимо стремиться к наименьшему объему такого кода в библиотеке или коде приложения. Небезопасность обнаруживается, выделяется и обозначается. Поэтому, если в вашем коде на Rust возникают segfaults, то вы нашли либо ошибку в компиляторе, либо ошибку в нескольких строках вашего небезопасного кода.

Это неидеальная система, но если вам нужна тройка - скорость, безопасность и многопоточность - это единственный возможный вариант.