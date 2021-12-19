# Abyss programming language

## Содержание
=============

## Описание
Abyss is a toy language interpreter, written in Go. It has C-style syntax, and is largely inspired by Ruby, Python, Perl and C#  
It support the normal control flow and functional programming.

## Использование
Чтобы запустить REPL, просто введите следующую команду:
```sh
~ >> abyss
Hello, <username>! This is the Abyss programming language!
Feel free to type in commands

>>
```

или запустить программу:
`abyss path/to/file`

## Тур по языку

### Типы данных

Abyss поддерживает базовые типы данных: `int`, `float`, `string`, `bool`, `array`, `hash`, `null`

```swift
s1 = "Hello" + " world!"
i  = 10
f  = 10.0
b  = true
a  = [1, "2", "three"]
t  = a[1 + 1 + 1]
h  = {"a": 1, "b": [1, true]}
n  = null
```

### Переменные
Чтобы создать переменную в Abyss используется ключевое слово `let`.

```swift
let a = [1, 2];
let b = "bob";
let c = true;
```

### Зарезервированные ключевые слова

Ключевые слова - это предопределенные зарезервированные идентификаторы, которые имеют особое значение для компилятора.  
Они не могут быть использованы в качестве идентификаторов. Ниже приведен список зарезервированных ключевых слов  

* fn
* let
* true false null
* if else
* while
* return

### Преобразование типов TODO

Вы можете использовать встроенные функции для преобразования типов переменных

* int()
* float()
* str()
* bool()

### Управляющая логика

Условия в Abyss являются выражениями, поэтому их можно присвоить переменным.
`let result = if (10 > 5) { true } else { false };`

* if / if-else
* while

```swift
// if-else
let a = 10;
let b = 5;
if (a == b) {
	print("a is equal to b");
} else {
	print("a is not equal to b");
}

// while
let i = 0
while (i < 10) {
	print(i);
	i += 1;
}
```

### Встроенные операторы

* Операторы сравнения: `<`, `<=`, `>`, `>=`, `==`, `!=`
* Унарные: `!`, `-`
* Арифметические: `+`, `-`, `*`, `/`
* Присваивания: `=`, `+=`, `-=`, `*=`, `/=`

### Функции

В Abyss функции - это функции первого класса.  
Это означает, что язык поддерживает передачу функций в качестве аргументов другим функциям,  
возвращает их в виде значений из других функций и присваивает их переменным.  
Также функции поддерживают *замыкание*!

```swift
// определить функцию
let add = fn(x) { x + 2; }; // return не обязателен
let a = [1, 2, add(1)];
print(a)

let complex = {
	"add": fn(x, y) { return fn(z) { x + y + z } }, // функция с замыканием!
	"sub": fn(x, y) { x - y },
	"other": [1, 2, 3, 4]
}
print(complex["add"](1, 2)(3))
print(complex["sub"](1, 2))
print(complex["other"][2])
```

#### Map

```swift
let map = fn(arr, f) {
	let iter = fn(arr, accumulated) {
		if (len(arr) == 0) {
			accumulated
		} else {
			iter(tail(arr), push(accumulated, f(first(arr))));
		}
	};

	iter(arr, []);
};

let a = [1, 2, 3, 4];
let double = fn(x) { x * 2 };
print(map(a, double));
```

Результат: `[2, 4, 6, 8]`

#### Reduce

```swift
let reduce = fn(arr, initial, f) {
	let iter = fn(arr, result) {
		if (len(arr) == 0) {
			result
		} else {
			iter(tail(arr), f(result, first(arr)));
		}
	};

	iter(arr, initial);
};
let sum = fn(arr) {
	reduce(arr, 0, fn(initial, el) { initial + el });
};

print(sum([1, 2, 3, 4, 5]));
```

Результат: `15`

### Встроенные функции
* `len(arr)` | `len(str)` -> вычисляет длину
* `first(arr)`			  -> возвращает первый элемент массива
* `tail(arr)`		      -> возвращает массив, без первого элемента
* `last(arr)`			  -> возвращает последний элемент массива
* `push(arr, value)`	  -> добавляет элемент в массив
* `array(size, value)`	  -> создает массив с указанным размером, где каждый элемент равен **value**
* `range(a, b)`			  -> возвращает массив с элементами от **a**, до **b**
* `print(value...)`       -> выводит на экран все аргументы
* `abs(value)`			  -> возвращает абсолютное значение числа
* `pow(base, exp)`		  -> возводит **base** в степень **exp**. Возвращает **float**
* `random()`			  -> рандомное вещественное число в интервале `[0.0, 1.0]`
* `sqrt(value)`           -> квадратный корень числа

