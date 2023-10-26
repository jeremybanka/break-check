# break-check
Semantic versioning tooling. Versatile testing function for finding breaking changes.

This work is the result of AI-based Test-Driven Development. The conversation started like so:

**me**
``` 
I have an idea for a npm package written in go called "break-check" which helps identify when breaking changes have been made in your code. The use case is,
- I'm on a feature branch, having made some changes.
- I run break-check, and my current code is tested using a subset of the tests that were present at the time of the last release—the subset should be only the tests covering my code's public API.
- I get one failed test, indicating I have introduced breaking changes into my code.
- Assuming my tests correctly cover my package's public interface, this tells me that I have made a breaking change.
Key considerations:
- The consumer should be able to define a search pattern for the subset of "public API tests", passed via CLI arg.
- The consumer should be able to define what the "test command" is.
- The tests that run should be the ones from the most recent tag in git history.
Can you walk me through implementing this from scratch?
```
**chatGPT**
```
Certainly! This is a complex, yet fascinating task. Implementing "break-check" from scratch will involve several steps, particularly:

Setting up the Go environment and initializing the project
Parsing CLI arguments and configuration
Identifying the latest Git tag
Fetching tests related to the last release using Git
Filtering tests based on the "public API" search pattern
Running the tests using the provided test command
Displaying the results
Let's walk through the process...
```
