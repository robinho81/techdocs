# techdocs

Tools to help build, publish and host a set of HTML documents for use as help files within an application e.g. tooltips and hints.

Pages can be generated from markdown documents in a git repository.
Hints or tooltips can be generated from text files. 

Both these get published against a particular version number for an application.

This golang project will allow the user to crawl directories for markdown documents and convert these to HTML.

## Blackfriday Markdown processor

This project uses the Blackfriday Markdown Processor

Available at https://github.com/russross/blackfriday

Blackfriday is distributed under the Simplified BSD License:

> Copyright © 2011 Russ Ross
> All rights reserved.
>
> Redistribution and use in source and binary forms, with or without
> modification, are permitted provided that the following conditions
> are met:
>
> 1.  Redistributions of source code must retain the above copyright
>     notice, this list of conditions and the following disclaimer.
>
> 2.  Redistributions in binary form must reproduce the above
>     copyright notice, this list of conditions and the following
>     disclaimer in the documentation and/or other materials provided with
>     the distribution.
>
> THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
> "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
> LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
> FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE
> COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
> INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
> BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
> LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
> CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
> LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN
> ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
> POSSIBILITY OF SUCH DAMAGE.
