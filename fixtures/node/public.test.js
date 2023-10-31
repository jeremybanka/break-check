const { describe, it } = require('node:test');

const assert = require('assert').strict;
const MyClass = require('./src.js');

describe('MyClass Public API Test', function() {
  it('should test publicMethod', function() {
    const obj = new MyClass();
    assert.strictEqual(obj.publicMethod(), "publicMethodOutput");
  });
});