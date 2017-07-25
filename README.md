Pontifex
========

Golang solitaire cypher from Cryptonomincon designed by Bruce Schneier 

More info https://www.schneier.com/academic/solitaire/

Quick start
-----------

1. Generate a sample key from the default deck key...

`go run main.go -op dkey`

this will make a base64 encoded key in a file called `shared.key`

2. or use and/or edit `deck.txt` to make a new key deck (it is not validated in any way)
generate a key from this deck:

`go run main.go -op key -k my.key -i deck.txt`


3. encrypt some text with the `shared.key`, for example encrypt the example `in.txt`

`go run main.go -op enc -k shared.key -i in.txt -o enc.txt`

cat the encoded text:

`cat enc.txt`

```
UPLOP IMCPO GIGJL LGSAX
GGGBS GXXLX
```

4. decrypt enc.txt with the `shared.key`:

`go run main.go -op dec -k shared.key -i enc.txt -o denc.txt`


Thats all there is to it. 

Testing
-------

The benchmark shows that we can generate 5.6e7 key stream values in 20 seconds, so should be pretty useful for randomness analysis (see Schneiers page for links regarding this)

License 
-------
This code is explicitly in the "public domain" and can be used, reused, and abused as desired.


THIS SOFTWARE IS PROVIDED "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN 
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS 
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.