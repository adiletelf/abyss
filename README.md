# Monkey

Statements:  
1. Let statement: `let x = 10;`  
2. Return statement: `return x;`  

Everything else is an expression.


### If-Expression:  
`let result = if (10 > 5) { true } else { false };`

### Map
```
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
map(a, double);

[2, 4, 6, 8]
```

### Reduce
```
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
sum([1, 2, 3, 4, 5]);

15
```
