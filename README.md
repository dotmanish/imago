imago-crop : image cropping tool written in Go
==============================================

I wrote this quick tool because I wanted a simple and narrow tool to crop images via command line. Using something all-in-one like ImageMagick or GraphicsMagick wasn't an option, and I was itching to write it myself. Hence this is the result. Hope it helps somebody learn or becomes a useful command line tool for someone's purpose.

Usage
=====

Crop the image test-input.jpg 100 px from top and 200 px from left and write the output to test-output.jpg

    imago-crop -i test-input.jpg -o test-output.jpg -top 100 -left 200

Crop the image test-input.jpg 50 px from right and write the output to test-output.png

    imago-crop -i test-input.jpg -o test-output.png -outformat png -right 50

Crop the image test-input.jpg 20% from top, 10% from bottom, and write the output to test-output.png

    imago-crop -i test-input.jpg -o test-output.png -outformat png -top 20% -bottom 10%

Crop the image test-input.jpg 25% px from left, ensuring at least 300px width,
and write the output to test-output.jpg

    imago-crop -i test-input.jpg -o test-output.jpg -left 25% -minwidth 300


Install
=======

You would want to do

    go get github.com/dotmanish/imago/imago-crop

to download and install the *imago-crop* binary.
This requires you to have configured GOPATH variable correctly in your
environment.


License
=======

Use of this source code is governed by a BSD (3-Clause) License.

Copyright 2013 Manish Malik (manishmalik.name)

All rights reserved.
    
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

    * Redistributions of source code must retain the above copyright notice,
      this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright notice,
      this list of conditions and the following disclaimer in the documentation
      and/or other materials provided with the distribution.
    * Neither the name of this program/product nor the names of its contributors may
      be used to endorse or promote products derived from this software without
      specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
