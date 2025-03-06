# Yappanese

Yappanese is a language interpreter using golang as its base language. The project is inspired by the "Writing an Interpreter In Go" book written by Thorsten Ball.
I used the book at the based for my Yappanese interpreter and adding more variables and functions that I often use while coding.
This is not quite a complete interpreter since there are still many things that I want to implement into it
Yappanese interpreter will read the file from top to bottom so unless you had not already declare your variables or function, it will give you an error.

## Requirement
You will need to have [golang](https://go.dev/doc/install) install in your machine.
I am currently running ver 1.22 so anything later than that should be working fine.

## Installation
You can repo this git with:
```bash
git clone https://github.com/ctrl-MizteRy/Yappanese.git
```

Or download it as a zip file in the GitHub page

## Usage
As you can see, there is a weird test.txt file sitting around in the repo, that would be where you write your code into.
Another option of testing out the language is to go into the main.go file and uncommented the import and the part of main, then comment out the other import and stuffs in main.

If you wish to run the test.txt file, do:
```bash
go run main.go test.txt
```
for any .txt file would be sufficient

To check out each input of the language and run the other option, just use:
```bash
go run .
```

## Declairation for variable

In Yappanese, you will be using the word `propose` to declare a variable.\
`propose a = 10;`\
You can also declare a null variable for later assign.\
`propose a;`\
`a = 10; #this is okay since a is already been declare`\
`a = "hello" # this will raise an error since a is now an integer obj and cannot be redeclare as a string obj`\

## Variables
Currently, Yappanese has **int64**, **float64**, **boolean**, **array**, and **hashmap**

### Int
You can basic do any simple arithmetic with Yappanese int var \
`
prosose a = 4;
` \
`
a = a + 5 # a will be 9
`\
`
a = a - 5 # a will be -1
`\
`
a = a * 3 # a will be 12
`\
`
a = a / 0 # this will raise an error, because you know why
`\
There is also the power and modulo:\
`a**2 = 16`\
`a%2 = 0`

### Float
Float will have the same arithmetic like int; however, if you do any arithmetic between int and float, the result will automatically convert into a float.
\
`
propose a = 4;
propose b = 1.0;
`\
`
a = a + b # a = 5.0
`\
`
a = a - b # a = 3.0
`\
`
a = a * b # a = 4.0
`\
`
a = a / b # a = 4.0
`\
Power and Modulo will be working the same in Float just as in Int;

### Boolean
Beside `true` and `false`, you will also have `nocap` and `cap`, which is pretty much the equivilant of the other 2.

### String
Currently Yappanese only support string and not character.
Declaration for string would be the same as how you do it in Go
`propose a = "hello"`\
String is currently support 2 type of arithmatic '+' and '*'
`"hello" + " world" = "hello world"`\
`"hello" * 3 = "hellohellohello"`

### Array
Yappanese is currently support nested array; however, array variables need to be from the same type. Even between float and int, this may change in the future tho.

`propose a = [1,2,3,4]`\
You can access each element using standard indexing.\
`propose b = a[3] # b = 4`

### HashMap
You can declare hashmap as:\
`propose a = {"hello": 1, "hi": 2, "aaaaaaa": -3};`\
As accessing it with:\
`propose b = a["aaaaaaa"] # b = -3`


## Function
In Yappanese, you there are 2 ways of declaring a function.

The standard way:\
`func name(param...) { whatever is in here}`\
You can later call the function with: `name(param)`\
\
The other way of doing it is:\
`propose a = func (param...){}`\
Then you can call the function later with:\
`a(param)`

## If Else Statements
Just like any other language, Yappanese support if else if and else.
The syntax for it would be:\
`if: perhaps`\
`else if: perchance`\
`else: otherwise`\
\
So if you want to write a if-else statement, it would look like this:\
```
perhaps (foo == nocap){
somthing in here
} perchance ( foo > bar) {
something else in there
} otherwise {
yap("hello world")
}
```
Yappanese is also support ternary expression for whoever wants to use it!
The syntax for this will be:\
`propose a = cap;`\
`propose b = (a)? 10 : 5 # b = 5`

## Loop
In Yappanese, the for loop will be be taking in a minimum of 1 parameter.
So you can use the for loop as a while loop just like in Go.

But you can also do something like this:\
`propose a = 3;`\
`for (a < 30, ++a){}`

## Builtin Functions
There are a couple of builtin function in Yappanese

### Len
This pretty much gonna give you the length of the obj, you can use it on string, array, and hashmap.
The syntax for len would be: `len(arr)`\

### Append
You can use this to add more varible into your array.
`propose a = [1,2,3,4]`\
`a = append(a, 5) # a = [1,2,3,4,5]`\
But you can also append to your array like this:\
`a = append(a, [5,6,7]) # a = [1,2,3,4,5,6,7]`

### Scan
Scan will pretty much take in whatever you type from the CLI. scan will not take any param for now.
Syntax:\
`propose a = scan();`\

### Pop
The pop function will take in 1 or 2 params. Array will always need to be in one of those.
The second parameter will be the index of where you want to pop the element out of the array.

Syntax:\
`propose a = [1,2,3,4,5]`\
With 1 param:\
`propose b = pop(a) # b = 5 and a = [1,2,3,4]`\
With 2 params:\
`propose b = pop(a,0) # b = 1, a = [2,3,4,5]`

### Keys
This function will return the keys of the hashmap in a form of array.
Syntax:\
`propose a = {"hello": 2, "hi": 1, "aaa": -5};`\
`propose b = keys(a) # b = ["hello", "hi", "aaa"]`\

### Values
Similar to keys, this function will return the value of each hashmap in the form of an array

### Yap
Last builtin function for Yappanese is of course yap, this is pretty much a printf function like in C.
However, unlike printf in C, you can add multiple variable into it and it will yap each element with a " " in between.\
Syntax:\
`yap("hello", "world") # you will have 'hello world' in your CLI`

## Contributing
I mean this is just a fun project that I write to learn Go and also how interpreter work.
If you want to do PR, I would not stop you but please also adding the proper testing file.

## Goals
The main goal for this moment to to write a complier for Yappanese.
