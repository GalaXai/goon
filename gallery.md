# example of filter use

<table>
<tr>
  <th>Original</th>
  <th>Downsampled</th>
  <th>Desaturated</th>
</tr>
<tr>
  <td><img src="static/test.png" style="width: 200px; height: 200px" alt="Original"></td>
  <td><img src="static/downsapled.png" style="width: 200px; height: 200px" alt="Downsampled"></td>
  <td><img src="static/desaturated.png" style="width: 200px; height: 200px" alt="Desaturated"></td>
</tr>
<tr>
  <th>diffOfGaussians</th>
  <th>sobelFilter</th>
  <th>sobelFilter Gradients</th>
</tr>
<tr>
  <td><img src="static/gaussDiff.png" style="width: 200px; height: 200px" alt="difference of gaussians"></td>
  <td><img src="static/sobelFilter.png" style="width: 200px; height: 200px" alt="sobel filter"></td>
  <td><img src="static/gradientMatrix.png" style="width: 200px; height: 200px" alt="sobelFilter Gradients"></td>
</tr>
</table>

# ascii
> also inside bash code for colors
```bash
?????????????????????????????????????????????????????????
?????????????????????????????????????????????????????????
?????????????????????????????????????????????????????????
????????????????????????????OOOOOOO??????????????????????
???????????????????????????OPPOPPoPPOO???????????????????
??????????????????????????PcccoPPo;oPOPO?????????????????
??????????????????OO??OPocccccooocccoooPO????????????????
??????????????Ooccc;;ocoocccccccc;P?;PPPOO???????????????
????????????Oocccc;;ccoc;;cc;cccc;P@?;PPPOO??????????????
???????????Occcc;;;ccccc;;;cccccc;P▓▓?coPOOO?OO??????????
????????O?Occcc;;,;ccccc;OP,;cccc;P▓▓▓OcPPPOOOO?O????????
????????cPoccc;,,;;cccc;,O@O,cccc;P▓▓▓▓?;oPOOOOOOOO??????
?????????Pccc;;,,;;cc;c;,O▓@O;ccc;P▓▓▓▓@?;oPOOPPO?OPO????
?????????ccc;;;,@;,cc;;;,O▓▓▓P;cc;P▓▓▓▓▓▓?coPOPPOO?OP????
??????O?P;;c;;,,@?c;cc;c,O▓▓▓▓O;c;P▓▓▓▓Poo;coPPoPO?O?O???
??????cc;;;;;;;,@▓?c;;;c;O▓▓▓▓▓O;;P▓@o▓;;;;coPPPooO??O???
??????cP,;;;;;;,@▓▓@;;;c;O▓▓▓▓▓@P,P@;;P@,cccoPPPoco???O??
??????c?O,;;;;;,@▓▓▓@c;c;O▓▓▓▓?OOoo;c;o@c;ccoPPPPocP??O??
?????O,?@?;;;;;,@▓▓▓▓@c;;O▓@O▓c,,,,ccc;?@,cccoOoP;;cO?O??
????Oc,?▓▓?;;;;,@▓▓▓▓▓?c,O@ocOO;ccccccc;;;cccoPPoo;;P?O??
????c;,?▓▓▓O;;;,@▓▓▓▓▓▓?;Oc;;o▓,cccccccc;ccccooooc;;oO???
???Pc;,?▓▓▓▓?,;,@▓@@@;,,,,ccc;O?,cccc;cccc;cccocc;;;cP???
???ccc,?▓▓▓▓@?,,@@oo@;;cccccc;oP,cccc;;cc;;cc;c;;;;;cO???
??Occc,?▓▓▓▓▓@O,@o;,?@,ccccccc;;cccccc;;c;;;ccc;;;;;oP???
??Occc,?▓▓▓▓OPP;o,;;o@c;ccccccccoccccc;;c;;ccco;;;;;oP???
??Occc,?▓?O▓;,,,;;;;,@O;ccccccccccccccc;ccccccc;,,;;oP???
???ccc,?@;;O?,;;;;;;;,;ccccccccccccccoo;occccoc;,;;;cO???
???oc;,?;c,P@,;;;;;c;;cccccccccccccoccocPccoco;;;,,;cO???
???Pc;,,;c;,?@,;c;;;;;c;c;ccccccccccoccoPooc;o,;;;,;o????
???Oc;;;;;;;oo,;;;;;c;cccccccccoccococcoooocc;,;,;,;PO???
????P;c;,;;;;,;;c;;;cccccccccccocooooccocOPc;,;,;;,;OOO??
?????oc;,,;;;;;;;;;;;cccccccccoooPoooccooPo;,,;;;;;cOOOO?
??????c;,,,;;;;;;;;cc;;;;ccccoocoPPooooPoPP;;;;;;;;oOOOO?
??????Oc,;;;;;;;;;;;;cccccocooccoPPPoPPPPPo,,;;;;;;POOOO?
???????Oc;;,,;,;;;;;ccccoooooocooPPOoPPPOOc;,,;;;;cOOOOO?
?????????P;;;;;;;;;;;;ccoooooooooPOOoPPOOOc;,;;;;;oOOOOO?
??????????Occcccccc;;coooooooPPPPPOOPOOOOPc;,;;cccOOOOOO?
??????????Pocccccc;;cPPPPPPPPPPPOPOOPOOOOPoc;;;;coOOOOO??
???????????Oocooo;ccPPOOOOOOOOOOOPOOPOOOOOPo;;,;cOOOOO???
?????????????POOOOPP?O?????OOOOOOOOOOOOOOOOPo;,;OOOO?????
?????????????OOOOOOO?????????OOOOOOOOOOOOOOOPc;cOO???????
?????????????OOOOOOO??????????OOOOOO?OOOOOOOOOooO????????
?????????????POOOOO????????????OOO??????OOOOOOPo?????????
????????????OOOOOOO???????????????????????OOOPPPO????????
????????????OOOOOOO?????????????????????????OPooo????????
?????????????OOPPOO????????????????????????OPc,,,O???????
????????????????O?????????????????????????Oc,,,,,o???????
????????????????O??????????????????????Poo;,,,,,,,???????
????????????????OOO???????????????????P,,,,,,,,,,,PO?????
?????????????????OOPOO??????????????Oc,,,,,,,,,,,,c??????
??????????????????O????????????OO??o;,,,,,,,,,,,,,,O?????
???????????????????OO?????????OOOP;,,,,,,,,,,,,,,,,c?????
????????????????????O????????OO??c,,,,,,,,,,,,,,,,,,o????
????????????????????????????OO???P,,,,,,,,,,,,,,,,,,,O???
?????????????????????O???????????O;,,,,,,,,,,,,,,,,,,,O??
?????????????????????????????????Oc,,,,,,,,,,,,,,,,,,,;P?
?????????????????????????????????Oc,,,,,,,,,,,,,,,,,,,,,c
```
## Cooking
```bash
                                                 \
                                                 \
                                                 \
                                                 \
                                                 \
                                                 \
                                                 \
                   \\\\\\\\\\\\                  \
                 \\\          \\\                \
               \\                \\              \
              \\                  \\\            \
            \\                      \\           \
           \\                        \\          \
           \                          \\         \
          \                            \         \
         \\                            \\        \
         \                              \        \
        \\                               \       \
        \                                \       \
        \                                \\      \
       \                                  \      \
       \                                  \      \
       \                                  \      \
       \                                  \      \
       \                                  \      \
       \                                  \      \
       \                                  \      \
       \                                  \      \
       \                                  \      \
       \\                                 \      \
        \                                \\      \
        \                                \       \
        \\                              \\       \
         \                              \        \
         \\                            \\        \
          \\                          \\         \
           \\                         \          \
            \\                       \           \
             \\                    \\            \
              \\                  \\             \
               \\\              \\\              \
                 \\\\        \\\\                \
                    \\\\\\\\\\                   \
                                                 \
                                                 \
                                                 \
                                                 \
                                                 \
                                                 \
```