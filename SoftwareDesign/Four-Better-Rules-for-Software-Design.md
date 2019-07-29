Улучшенные четыре правила дизайна ПО
David Bryant Copeland
https://naildrivin5.com/blog/2019/07/25/four-better-rules-for-software-design.html

Мартин Фаулер недавно создал твит с ссылкой на его [пост в блоге](https://martinfowler.com/bliki/BeckDesignRules.html) о четырех правилах простого дизайна от Кента Бека, которые, как я думаю, могут быть еще улучшены (и которые иногда могут отправить программиста по ложному пути):

Правила Кента, из книги [Extreme Programming Explained](https://www.amazon.com/gp/product/0201616416):

- Кент говорит: "Запускайте все тесты".
- Не дублируйте логику. Старайтесь избегать скрытых дубликатов, таких как параллельных иерархий классов.
- Все намерения, важные для программиста, должны быть явно видны.
- Код должен иметь наименьшее возможное количество классов и методов.

Согласно моему опыту, эти правила не совсем соответствуют нуждам дизайна ПО. Мои четыре правила хорошо спроектированной системы могли бы быть такими:

- она хорошо покрывается тестами и успешно проходит их.
- она не имеет абстракций, которые напрямую не нужны программе.
- она имеет однозначное поведение.
- она требует наименьшее количество концепций.

Для меня эти правила вытекают из того, что мы делаем с нашим ПО.

# Что мы *делаем* с нашим ПО?

Мы не можем говорить о дизайне ПО, не поговорив прежде о том, что мы намереваемся делать с этим ПО.

ПО пишется для решения проблемы. Программа выполняется и имеет поведение. Это поведение изучается, чтобы обеспечить правильность работы или обнаружить ошибки. ПО также часто изменяется для придания ему нового или другого поведения.

Поэтому любой подход к дизайну ПО должен быть сфокусирован на предсказании, изучении и понимании его поведения, чтобы сделать изменение этого поведения как можно проще.

Мы проверяем правильность поведения путем тестирования, и поэтому я согласен с Кентом, что первое и самое главное - хорошо спроектированное ПО должно проходить тесты. Я даже пойду дальше и настою на том, что ПО должно иметь тесты (т. е. быть хорошо покрытым тестами).

После того, как поведение было проверено, следующие три пункта обоих списков касаются понимания нашего ПО (и, следовательно, его поведения). Его список начинается с дублирования кода, которое действительно находится на своем месте. Однако по моему личному опыту, слишком сильный фокус на уменьшение дублирования кода имеет высокую цену. Для его устранения необходимо создать скрывающие его абстракции, и именно эти абстракции делают ПО сложным для понимания и изменения.

# Устранение дублирования кода требует абстракций, а абстракции приводят к сложности

"Don't Repeat Yourself" или DRY используется для оправдания спорных решений дизайна. Вы когда-нибудь видели подобный код?

'''ruby
ZERO = BigDecimal.new(0)
'''

Кроме того вы, наверное, видели что-то подобное:

'''java
public void call(Map payload, boolean async, int errorStrategy) {
  // ...
}
'''

Если вы видите методы или функции с флагами, boolean и т.д., то это, как правило, значит, что кто-то использовал принцип DRY при рефакторинге, но код не был *точно* таким же в обоих местах, поэтому полученный код должен быть достаточно гибким, чтобы вмешать оба поведения.

Такие обобщенные абстракции сложны для тестирования и понимания, так как они должны обрабатывать намного больше случаев, чем оригинальный (возможно, продублированный) код. Иными словами, абстракции поддерживают намного больше поведения, чем нужно для нормального функционирования системы. Таким образом, устранение дублирования кода может породить новое поведение, которое не требуется системе.

Поэтому *действительно* важно объединять некоторые типы поведения, однако, бывает сложно понять, какое поведение действительно дублируется. Часто куски кода выглядят похоже, но это происходит только по случайности.

Подумайте, насколько проще устранить дублирование кода, чем повторно вернуть его (например, после создания плохо продуманной абстракции). Поэтому нам нужно нужно задуматься об оставлении дублирования кода, если только мы не абсолютно уверены, что у нас есть лучший способ избавиться от него.

Создание абстракций должно заставить нас задуматься. Если в процессе устранения дублирующегося кода вы создаете очень гибкую обобщенную абстракцию, то, возможно, вы пошли по неверному пути.

Это приводит нас к следующему пункту - намерение против поведения.

# Намерение программиста бессмысленно - поведение значит все