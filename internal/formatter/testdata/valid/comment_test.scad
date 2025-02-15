/*
   Hello, this is a multiline comment.
   It describes some stuff about the file
*/

// sometimes we use single line comments
// to comment on something that takes
// longer than a line

c = [ [0,0], // first row
[1,1], // second row
[2,2] // third row
];

module foo(argument1) {
  // this is a for loop
  for(i = [0:10]) {
    doSomething();// do something here!

// ====
    // below is a call to anotherFunction
    // ====
    anotherFunction();

    echo(i); /* multiline commends at the end of line are not correctly handled */

    if(i % 2 == 0) {
      echo("i is even"); // single line comment with a preceding space
      echo("which is good");
    }
  }
}

// comment 1
module x() { // comment 2
  // comment 3
} // comment 4
// comment 5
