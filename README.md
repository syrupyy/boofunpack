# boofunpack
Cocos2d-x texture unpacker, primarily for Cookie Run: OvenBreak and Cookie Run for Kakao/LINE (though it likely works for other .plist format cocos2d-x games too).
## Installation
Grab an executable for your OS on the [Releases](https://github.com/syrupyy/boofunpack/releases) page. (64-bit builds only for now, sorry for the inconvenience. I'm sure most people will just be using the 64-bit Windows version, though...)
## Building
1. [Install Go if you haven't already](https://go.dev/)
2. Clone the repository and install the dependencies, usually with `go get ./...`
3. Build it with `go build` (I usually also add the argument `-ldflags="-s -w"` for smaller executables) or run it with `go run main.go`
## Known Issues
- JSON files are not supported. Cookie Run doesn't use them, but games that do will need that support to be added.
- Kakao/LINE cocos2d files aren't supported yet, but they will be soon!!!!!! Probably.
## License
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.