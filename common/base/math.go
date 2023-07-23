package base

import (
	"fmt"
	"math"
)

const (
	denominator int32 = 10000
)

// 三角函數值
var (
	sin0_00  int32 = 0
	sin0_10  int32 = 29
	sin0_20  int32 = 58
	sin0_30  int32 = 87
	sin0_40  int32 = 116
	sin0_50  int32 = 145
	sin1_00  int32 = 175
	sin1_10  int32 = 204
	sin1_20  int32 = 233
	sin1_30  int32 = 262
	sin1_40  int32 = 291
	sin1_50  int32 = 320
	sin2_00  int32 = 349
	sin2_10  int32 = 378
	sin2_20  int32 = 407
	sin2_30  int32 = 436
	sin2_40  int32 = 465
	sin2_50  int32 = 494
	sin3_00  int32 = 523
	sin3_10  int32 = 552
	sin3_20  int32 = 581
	sin3_30  int32 = 610
	sin3_40  int32 = 640
	sin3_50  int32 = 669
	sin4_00  int32 = 698
	sin4_10  int32 = 727
	sin4_20  int32 = 756
	sin4_30  int32 = 785
	sin4_40  int32 = 814
	sin4_50  int32 = 843
	sin5_00  int32 = 872
	sin5_10  int32 = 901
	sin5_20  int32 = 929
	sin5_30  int32 = 958
	sin5_40  int32 = 987
	sin5_50  int32 = 1016
	sin6_00  int32 = 1045
	sin6_10  int32 = 1074
	sin6_20  int32 = 1103
	sin6_30  int32 = 1132
	sin6_40  int32 = 1161
	sin6_50  int32 = 1190
	sin7_00  int32 = 1219
	sin7_10  int32 = 1248
	sin7_20  int32 = 1276
	sin7_30  int32 = 1305
	sin7_40  int32 = 1334
	sin7_50  int32 = 1363
	sin8_00  int32 = 1392
	sin8_10  int32 = 1421
	sin8_20  int32 = 1449
	sin8_30  int32 = 1478
	sin8_40  int32 = 1507
	sin8_50  int32 = 1536
	sin9_00  int32 = 1564
	sin9_10  int32 = 1593
	sin9_20  int32 = 1622
	sin9_30  int32 = 1650
	sin9_40  int32 = 1679
	sin9_50  int32 = 1708
	sin10_00 int32 = 1736
	sin10_10 int32 = 1765
	sin10_20 int32 = 1794
	sin10_30 int32 = 1822
	sin10_40 int32 = 1851
	sin10_50 int32 = 1880
	sin11_00 int32 = 1908
	sin11_10 int32 = 1937
	sin11_20 int32 = 1965
	sin11_30 int32 = 1994
	sin11_40 int32 = 2022
	sin11_50 int32 = 2051
	sin12_00 int32 = 2079
	sin12_10 int32 = 2108
	sin12_20 int32 = 2136
	sin12_30 int32 = 2164
	sin12_40 int32 = 2193
	sin12_50 int32 = 2221
	sin13_00 int32 = 2250
	sin13_10 int32 = 2278
	sin13_20 int32 = 2306
	sin13_30 int32 = 2334
	sin13_40 int32 = 2363
	sin13_50 int32 = 2391
	sin14_00 int32 = 2419
	sin14_10 int32 = 2447
	sin14_20 int32 = 2476
	sin14_30 int32 = 2504
	sin14_40 int32 = 2532
	sin14_50 int32 = 2560
	sin15_00 int32 = 2588
	sin15_10 int32 = 2616
	sin15_20 int32 = 2644
	sin15_30 int32 = 2672
	sin15_40 int32 = 2700
	sin15_50 int32 = 2728
	sin16_00 int32 = 2756
	sin16_10 int32 = 2784
	sin16_20 int32 = 2812
	sin16_30 int32 = 2840
	sin16_40 int32 = 2868
	sin16_50 int32 = 2896
	sin17_00 int32 = 2924
	sin17_10 int32 = 2952
	sin17_20 int32 = 2979
	sin17_30 int32 = 3007
	sin17_40 int32 = 3035
	sin17_50 int32 = 3062
	sin18_00 int32 = 3090
	sin18_10 int32 = 3118
	sin18_20 int32 = 3145
	sin18_30 int32 = 3173
	sin18_40 int32 = 3201
	sin18_50 int32 = 3228
	sin19_00 int32 = 3256
	sin19_10 int32 = 3283
	sin19_20 int32 = 3311
	sin19_30 int32 = 3338
	sin19_40 int32 = 3365
	sin19_50 int32 = 3393
	sin20_00 int32 = 3420
	sin20_10 int32 = 3448
	sin20_20 int32 = 3475
	sin20_30 int32 = 3502
	sin20_40 int32 = 3529
	sin20_50 int32 = 3557
	sin21_00 int32 = 3584
	sin21_10 int32 = 3611
	sin21_20 int32 = 3638
	sin21_30 int32 = 3665
	sin21_40 int32 = 3692
	sin21_50 int32 = 3719
	sin22_00 int32 = 3746
	sin22_10 int32 = 3773
	sin22_20 int32 = 3800
	sin22_30 int32 = 3827
	sin22_40 int32 = 3854
	sin22_50 int32 = 3881
	sin23_00 int32 = 3907
	sin23_10 int32 = 3934
	sin23_20 int32 = 3961
	sin23_30 int32 = 3987
	sin23_40 int32 = 4014
	sin23_50 int32 = 4041
	sin24_00 int32 = 4067
	sin24_10 int32 = 4094
	sin24_20 int32 = 4120
	sin24_30 int32 = 4147
	sin24_40 int32 = 4173
	sin24_50 int32 = 4200
	sin25_00 int32 = 4226
	sin25_10 int32 = 4253
	sin25_20 int32 = 4279
	sin25_30 int32 = 4305
	sin25_40 int32 = 4331
	sin25_50 int32 = 4358
	sin26_00 int32 = 4384
	sin26_10 int32 = 4410
	sin26_20 int32 = 4436
	sin26_30 int32 = 4462
	sin26_40 int32 = 4488
	sin26_50 int32 = 4514
	sin27_00 int32 = 4540
	sin27_10 int32 = 4566
	sin27_20 int32 = 4592
	sin27_30 int32 = 4617
	sin27_40 int32 = 4643
	sin27_50 int32 = 4669
	sin28_00 int32 = 4695
	sin28_10 int32 = 4720
	sin28_20 int32 = 4746
	sin28_30 int32 = 4772
	sin28_40 int32 = 4797
	sin28_50 int32 = 4823
	sin29_00 int32 = 4848
	sin29_10 int32 = 4874
	sin29_20 int32 = 4899
	sin29_30 int32 = 4924
	sin29_40 int32 = 4950
	sin29_50 int32 = 4975
	sin30_00 int32 = 5000
	sin30_10 int32 = 5025
	sin30_20 int32 = 5050
	sin30_30 int32 = 5075
	sin30_40 int32 = 5100
	sin30_50 int32 = 5125
	sin31_00 int32 = 5150
	sin31_10 int32 = 5175
	sin31_20 int32 = 5200
	sin31_30 int32 = 5225
	sin31_40 int32 = 5250
	sin31_50 int32 = 5275
	sin32_00 int32 = 5299
	sin32_10 int32 = 5324
	sin32_20 int32 = 5348
	sin32_30 int32 = 5373
	sin32_40 int32 = 5398
	sin32_50 int32 = 5422
	sin33_00 int32 = 5446
	sin33_10 int32 = 5471
	sin33_20 int32 = 5495
	sin33_30 int32 = 5519
	sin33_40 int32 = 5544
	sin33_50 int32 = 5568
	sin34_00 int32 = 5592
	sin34_10 int32 = 5616
	sin34_20 int32 = 5640
	sin34_30 int32 = 5664
	sin34_40 int32 = 5688
	sin34_50 int32 = 5712
	sin35_00 int32 = 5736
	sin35_10 int32 = 5760
	sin35_20 int32 = 5783
	sin35_30 int32 = 5807
	sin35_40 int32 = 5831
	sin35_50 int32 = 5854
	sin36_00 int32 = 5878
	sin36_10 int32 = 5901
	sin36_20 int32 = 5925
	sin36_30 int32 = 5948
	sin36_40 int32 = 5972
	sin36_50 int32 = 5995
	sin37_00 int32 = 6018
	sin37_10 int32 = 6041
	sin37_20 int32 = 6065
	sin37_30 int32 = 6088
	sin37_40 int32 = 6111
	sin37_50 int32 = 6134
	sin38_00 int32 = 6157
	sin38_10 int32 = 6180
	sin38_20 int32 = 6202
	sin38_30 int32 = 6225
	sin38_40 int32 = 6248
	sin38_50 int32 = 6271
	sin39_00 int32 = 6293
	sin39_10 int32 = 6316
	sin39_20 int32 = 6338
	sin39_30 int32 = 6361
	sin39_40 int32 = 6383
	sin39_50 int32 = 6406
	sin40_00 int32 = 6428
	sin40_10 int32 = 6450
	sin40_20 int32 = 6472
	sin40_30 int32 = 6494
	sin40_40 int32 = 6517
	sin40_50 int32 = 6539
	sin41_00 int32 = 6561
	sin41_10 int32 = 6583
	sin41_20 int32 = 6604
	sin41_30 int32 = 6626
	sin41_40 int32 = 6648
	sin41_50 int32 = 6670
	sin42_00 int32 = 6691
	sin42_10 int32 = 6713
	sin42_20 int32 = 6734
	sin42_30 int32 = 6756
	sin42_40 int32 = 6777
	sin42_50 int32 = 6799
	sin43_00 int32 = 6820
	sin43_10 int32 = 6841
	sin43_20 int32 = 6862
	sin43_30 int32 = 6884
	sin43_40 int32 = 6905
	sin43_50 int32 = 6926
	sin44_00 int32 = 6947
	sin44_10 int32 = 6967
	sin44_20 int32 = 6988
	sin44_30 int32 = 7009
	sin44_40 int32 = 7030
	sin44_50 int32 = 7050
	sin45_00 int32 = 7071
	sin45_10 int32 = 7092
	sin45_20 int32 = 7112
	sin45_30 int32 = 7133
	sin45_40 int32 = 7153
	sin45_50 int32 = 7173
	sin46_00 int32 = 7193
	sin46_10 int32 = 7214
	sin46_20 int32 = 7234
	sin46_30 int32 = 7254
	sin46_40 int32 = 7274
	sin46_50 int32 = 7294
	sin47_00 int32 = 7314
	sin47_10 int32 = 7333
	sin47_20 int32 = 7353
	sin47_30 int32 = 7373
	sin47_40 int32 = 7392
	sin47_50 int32 = 7412
	sin48_00 int32 = 7431
	sin48_10 int32 = 7451
	sin48_20 int32 = 7470
	sin48_30 int32 = 7490
	sin48_40 int32 = 7509
	sin48_50 int32 = 7528
	sin49_00 int32 = 7547
	sin49_10 int32 = 7566
	sin49_20 int32 = 7585
	sin49_30 int32 = 7604
	sin49_40 int32 = 7623
	sin49_50 int32 = 7642
	sin50_00 int32 = 7660
	sin50_10 int32 = 7679
	sin50_20 int32 = 7698
	sin50_30 int32 = 7716
	sin50_40 int32 = 7735
	sin50_50 int32 = 7753
	sin51_00 int32 = 7771
	sin51_10 int32 = 7790
	sin51_20 int32 = 7808
	sin51_30 int32 = 7826
	sin51_40 int32 = 7844
	sin51_50 int32 = 7862
	sin52_00 int32 = 7880
	sin52_10 int32 = 7898
	sin52_20 int32 = 7916
	sin52_30 int32 = 7934
	sin52_40 int32 = 7951
	sin52_50 int32 = 7969
	sin53_00 int32 = 7986
	sin53_10 int32 = 8004
	sin53_20 int32 = 8021
	sin53_30 int32 = 8039
	sin53_40 int32 = 8056
	sin53_50 int32 = 8073
	sin54_00 int32 = 8090
	sin54_10 int32 = 8107
	sin54_20 int32 = 8124
	sin54_30 int32 = 8141
	sin54_40 int32 = 8158
	sin54_50 int32 = 8175
	sin55_00 int32 = 8192
	sin55_10 int32 = 8208
	sin55_20 int32 = 8225
	sin55_30 int32 = 8241
	sin55_40 int32 = 8258
	sin55_50 int32 = 8274
	sin56_00 int32 = 8290
	sin56_10 int32 = 8307
	sin56_20 int32 = 8323
	sin56_30 int32 = 8339
	sin56_40 int32 = 8355
	sin56_50 int32 = 8371
	sin57_00 int32 = 8387
	sin57_10 int32 = 8403
	sin57_20 int32 = 8418
	sin57_30 int32 = 8434
	sin57_40 int32 = 8450
	sin57_50 int32 = 8465
	sin58_00 int32 = 8480
	sin58_10 int32 = 8496
	sin58_20 int32 = 8511
	sin58_30 int32 = 8526
	sin58_40 int32 = 8542
	sin58_50 int32 = 8557
	sin59_00 int32 = 8572
	sin59_10 int32 = 8587
	sin59_20 int32 = 8601
	sin59_30 int32 = 8616
	sin59_40 int32 = 8631
	sin59_50 int32 = 8646
	sin60_00 int32 = 8660
	sin60_10 int32 = 8675
	sin60_20 int32 = 8689
	sin60_30 int32 = 8704
	sin60_40 int32 = 8718
	sin60_50 int32 = 8732
	sin61_00 int32 = 8746
	sin61_10 int32 = 8760
	sin61_20 int32 = 8774
	sin61_30 int32 = 8788
	sin61_40 int32 = 8802
	sin61_50 int32 = 8816
	sin62_00 int32 = 8829
	sin62_10 int32 = 8843
	sin62_20 int32 = 8857
	sin62_30 int32 = 8870
	sin62_40 int32 = 8884
	sin62_50 int32 = 8897
	sin63_00 int32 = 8910
	sin63_10 int32 = 8923
	sin63_20 int32 = 8936
	sin63_30 int32 = 8949
	sin63_40 int32 = 8962
	sin63_50 int32 = 8975
	sin64_00 int32 = 8988
	sin64_10 int32 = 9001
	sin64_20 int32 = 9013
	sin64_30 int32 = 9026
	sin64_40 int32 = 9038
	sin64_50 int32 = 9051
	sin65_00 int32 = 9063
	sin65_10 int32 = 9075
	sin65_20 int32 = 9088
	sin65_30 int32 = 9100
	sin65_40 int32 = 9112
	sin65_50 int32 = 9124
	sin66_00 int32 = 9135
	sin66_10 int32 = 9147
	sin66_20 int32 = 9159
	sin66_30 int32 = 9171
	sin66_40 int32 = 9182
	sin66_50 int32 = 9194
	sin67_00 int32 = 9205
	sin67_10 int32 = 9216
	sin67_20 int32 = 9228
	sin67_30 int32 = 9239
	sin67_40 int32 = 9250
	sin67_50 int32 = 9261
	sin68_00 int32 = 9272
	sin68_10 int32 = 9283
	sin68_20 int32 = 9293
	sin68_30 int32 = 9304
	sin68_40 int32 = 9315
	sin68_50 int32 = 9325
	sin69_00 int32 = 9336
	sin69_10 int32 = 9346
	sin69_20 int32 = 9356
	sin69_30 int32 = 9367
	sin69_40 int32 = 9377
	sin69_50 int32 = 9387
	sin70_00 int32 = 9397
	sin70_10 int32 = 9407
	sin70_20 int32 = 9417
	sin70_30 int32 = 9426
	sin70_40 int32 = 9436
	sin70_50 int32 = 9446
	sin71_00 int32 = 9455
	sin71_10 int32 = 9465
	sin71_20 int32 = 9474
	sin71_30 int32 = 9483
	sin71_40 int32 = 9492
	sin71_50 int32 = 9502
	sin72_00 int32 = 9511
	sin72_10 int32 = 9520
	sin72_20 int32 = 9528
	sin72_30 int32 = 9537
	sin72_40 int32 = 9546
	sin72_50 int32 = 9555
	sin73_00 int32 = 9563
	sin73_10 int32 = 9572
	sin73_20 int32 = 9580
	sin73_30 int32 = 9588
	sin73_40 int32 = 9596
	sin73_50 int32 = 9605
	sin74_00 int32 = 9613
	sin74_10 int32 = 9621
	sin74_20 int32 = 9628
	sin74_30 int32 = 9636
	sin74_40 int32 = 9644
	sin74_50 int32 = 9652
	sin75_00 int32 = 9659
	sin75_10 int32 = 9667
	sin75_20 int32 = 9674
	sin75_30 int32 = 9681
	sin75_40 int32 = 9689
	sin75_50 int32 = 9696
	sin76_00 int32 = 9703
	sin76_10 int32 = 9710
	sin76_20 int32 = 9717
	sin76_30 int32 = 9724
	sin76_40 int32 = 9730
	sin76_50 int32 = 9737
	sin77_00 int32 = 9744
	sin77_10 int32 = 9750
	sin77_20 int32 = 9757
	sin77_30 int32 = 9763
	sin77_40 int32 = 9769
	sin77_50 int32 = 9775
	sin78_00 int32 = 9781
	sin78_10 int32 = 9787
	sin78_20 int32 = 9793
	sin78_30 int32 = 9799
	sin78_40 int32 = 9805
	sin78_50 int32 = 9811
	sin79_00 int32 = 9816
	sin79_10 int32 = 9822
	sin79_20 int32 = 9827
	sin79_30 int32 = 9833
	sin79_40 int32 = 9838
	sin79_50 int32 = 9843
	sin80_00 int32 = 9848
	sin80_10 int32 = 9853
	sin80_20 int32 = 9858
	sin80_30 int32 = 9863
	sin80_40 int32 = 9868
	sin80_50 int32 = 9872
	sin81_00 int32 = 9877
	sin81_10 int32 = 9881
	sin81_20 int32 = 9886
	sin81_30 int32 = 9890
	sin81_40 int32 = 9894
	sin81_50 int32 = 9899
	sin82_00 int32 = 9903
	sin82_10 int32 = 9907
	sin82_20 int32 = 9911
	sin82_30 int32 = 9914
	sin82_40 int32 = 9918
	sin82_50 int32 = 9922
	sin83_00 int32 = 9925
	sin83_10 int32 = 9929
	sin83_20 int32 = 9932
	sin83_30 int32 = 9936
	sin83_40 int32 = 9939
	sin83_50 int32 = 9942
	sin84_00 int32 = 9945
	sin84_10 int32 = 9948
	sin84_20 int32 = 9951
	sin84_30 int32 = 9954
	sin84_40 int32 = 9957
	sin84_50 int32 = 9959
	sin85_00 int32 = 9962
	sin85_10 int32 = 9964
	sin85_20 int32 = 9967
	sin85_30 int32 = 9969
	sin85_40 int32 = 9971
	sin85_50 int32 = 9974
	sin86_00 int32 = 9976
	sin86_10 int32 = 9978
	sin86_20 int32 = 9980
	sin86_30 int32 = 9981
	sin86_40 int32 = 9983
	sin86_50 int32 = 9985
	sin87_00 int32 = 9986
	sin87_10 int32 = 9988
	sin87_20 int32 = 9989
	sin87_30 int32 = 9990
	sin87_40 int32 = 9992
	sin87_50 int32 = 9993
	sin88_00 int32 = 9994
	sin88_10 int32 = 9995
	sin88_20 int32 = 9996
	sin88_30 int32 = 9997
	sin88_40 int32 = 9997
	sin88_50 int32 = 9998
	sin89_00 int32 = 9998
	sin89_10 int32 = 9999
	sin89_20 int32 = 9999
	sin89_30 int32 = 10000
	sin89_40 int32 = 10000
	sin89_50 int32 = 10000
	sin90_00 int32 = 10000
)

var sinval = [][6]int32{
	{sin0_00, sin0_10, sin0_20, sin0_30, sin0_40, sin0_50},
	{sin1_00, sin1_10, sin1_20, sin1_30, sin1_40, sin1_50},
	{sin2_00, sin2_10, sin2_20, sin2_30, sin2_40, sin2_50},
	{sin3_00, sin3_10, sin3_20, sin3_30, sin3_40, sin3_50},
	{sin4_00, sin4_10, sin4_20, sin4_30, sin4_40, sin4_50},
	{sin5_00, sin5_10, sin5_20, sin5_30, sin5_40, sin5_50},
	{sin6_00, sin6_10, sin6_20, sin6_30, sin6_40, sin6_50},
	{sin7_00, sin7_10, sin7_20, sin7_30, sin7_40, sin7_50},
	{sin8_00, sin8_10, sin8_20, sin8_30, sin8_40, sin8_50},
	{sin9_00, sin9_10, sin9_20, sin9_30, sin9_40, sin9_50},
	{sin10_00, sin10_10, sin10_20, sin10_30, sin10_40, sin10_50}, // 10
	{sin11_00, sin11_10, sin11_20, sin11_30, sin11_40, sin11_50},
	{sin12_00, sin12_10, sin12_20, sin12_30, sin12_40, sin12_50},
	{sin13_00, sin13_10, sin13_20, sin13_30, sin13_40, sin13_50},
	{sin14_00, sin14_10, sin14_20, sin14_30, sin14_40, sin14_50},
	{sin15_00, sin15_10, sin15_20, sin15_30, sin15_40, sin15_50},
	{sin16_00, sin16_10, sin16_20, sin16_30, sin16_40, sin16_50},
	{sin17_00, sin17_10, sin17_20, sin17_30, sin17_40, sin17_50},
	{sin18_00, sin18_10, sin18_20, sin18_30, sin18_40, sin18_50},
	{sin19_00, sin19_10, sin19_20, sin19_30, sin19_40, sin19_50},
	{sin20_00, sin20_10, sin20_20, sin20_30, sin20_40, sin20_50}, // 20
	{sin21_00, sin21_10, sin21_20, sin21_30, sin21_40, sin21_50},
	{sin22_00, sin22_10, sin22_20, sin22_30, sin22_40, sin22_50},
	{sin23_00, sin23_10, sin23_20, sin23_30, sin23_40, sin23_50},
	{sin24_00, sin24_10, sin24_20, sin24_30, sin24_40, sin24_50},
	{sin25_00, sin25_10, sin25_20, sin25_30, sin25_40, sin25_50},
	{sin26_00, sin26_10, sin26_20, sin26_30, sin26_40, sin26_50},
	{sin27_00, sin27_10, sin27_20, sin27_30, sin27_40, sin27_50},
	{sin28_00, sin28_10, sin28_20, sin28_30, sin28_40, sin28_50},
	{sin29_00, sin29_10, sin29_20, sin29_30, sin29_40, sin29_50},
	{sin30_00, sin30_10, sin30_20, sin30_30, sin30_40, sin30_50}, // 30
	{sin31_00, sin31_10, sin31_20, sin31_30, sin31_40, sin31_50},
	{sin32_00, sin32_10, sin32_20, sin32_30, sin32_40, sin32_50},
	{sin33_00, sin33_10, sin33_20, sin33_30, sin33_40, sin33_50},
	{sin34_00, sin34_10, sin34_20, sin34_30, sin34_40, sin34_50},
	{sin35_00, sin35_10, sin35_20, sin35_30, sin35_40, sin35_50},
	{sin36_00, sin36_10, sin36_20, sin36_30, sin36_40, sin36_50},
	{sin37_00, sin37_10, sin37_20, sin37_30, sin37_40, sin37_50},
	{sin38_00, sin38_10, sin38_20, sin38_30, sin38_40, sin38_50},
	{sin39_00, sin39_10, sin39_20, sin39_30, sin39_40, sin39_50},
	{sin40_00, sin40_10, sin40_20, sin40_30, sin40_40, sin40_50}, // 40
	{sin41_00, sin41_10, sin41_20, sin41_30, sin41_40, sin41_50},
	{sin42_00, sin42_10, sin42_20, sin42_30, sin42_40, sin42_50},
	{sin43_00, sin43_10, sin43_20, sin43_30, sin43_40, sin43_50},
	{sin44_00, sin44_10, sin44_20, sin44_30, sin44_40, sin44_50},
	{sin45_00, sin45_10, sin45_20, sin45_30, sin45_40, sin45_50},
	{sin46_00, sin46_10, sin46_20, sin46_30, sin46_40, sin46_50},
	{sin47_00, sin47_10, sin47_20, sin47_30, sin47_40, sin47_50},
	{sin48_00, sin48_10, sin48_20, sin48_30, sin48_40, sin48_50},
	{sin49_00, sin49_10, sin49_20, sin49_30, sin49_40, sin49_50},
	{sin50_00, sin50_10, sin50_20, sin50_30, sin50_40, sin50_50}, // 50
	{sin51_00, sin51_10, sin51_20, sin51_30, sin51_40, sin51_50},
	{sin52_00, sin52_10, sin52_20, sin52_30, sin52_40, sin52_50},
	{sin53_00, sin53_10, sin53_20, sin53_30, sin53_40, sin53_50},
	{sin54_00, sin54_10, sin54_20, sin54_30, sin54_40, sin54_50},
	{sin55_00, sin55_10, sin55_20, sin55_30, sin55_40, sin55_50},
	{sin56_00, sin56_10, sin56_20, sin56_30, sin56_40, sin56_50},
	{sin57_00, sin57_10, sin57_20, sin57_30, sin57_40, sin57_50},
	{sin58_00, sin58_10, sin58_20, sin58_30, sin58_40, sin58_50},
	{sin59_00, sin59_10, sin59_20, sin59_30, sin59_40, sin59_50},
	{sin60_00, sin60_10, sin60_20, sin60_30, sin60_40, sin60_50}, // 60
	{sin61_00, sin61_10, sin61_20, sin61_30, sin61_40, sin61_50},
	{sin62_00, sin62_10, sin62_20, sin62_30, sin62_40, sin62_50},
	{sin63_00, sin63_10, sin63_20, sin63_30, sin63_40, sin63_50},
	{sin64_00, sin64_10, sin64_20, sin64_30, sin64_40, sin64_50},
	{sin65_00, sin65_10, sin65_20, sin65_30, sin65_40, sin65_50},
	{sin66_00, sin66_10, sin66_20, sin66_30, sin66_40, sin66_50},
	{sin67_00, sin67_10, sin67_20, sin67_30, sin67_40, sin67_50},
	{sin68_00, sin68_10, sin68_20, sin68_30, sin68_40, sin68_50},
	{sin69_00, sin69_10, sin69_20, sin69_30, sin69_40, sin69_50},
	{sin70_00, sin70_10, sin70_20, sin70_30, sin70_40, sin70_50}, // 70
	{sin71_00, sin71_10, sin71_20, sin71_30, sin71_40, sin71_50},
	{sin72_00, sin72_10, sin72_20, sin72_30, sin72_40, sin72_50},
	{sin73_00, sin73_10, sin73_20, sin73_30, sin73_40, sin73_50},
	{sin74_00, sin74_10, sin74_20, sin74_30, sin74_40, sin74_50},
	{sin75_00, sin75_10, sin75_20, sin75_30, sin75_40, sin75_50},
	{sin76_00, sin76_10, sin76_20, sin76_30, sin76_40, sin76_50},
	{sin77_00, sin77_10, sin77_20, sin77_30, sin77_40, sin77_50},
	{sin78_00, sin78_10, sin78_20, sin78_30, sin78_40, sin78_50},
	{sin79_00, sin79_10, sin79_20, sin79_30, sin79_40, sin79_50},
	{sin80_00, sin80_10, sin80_20, sin80_30, sin80_40, sin80_50}, // 80
	{sin81_00, sin81_10, sin81_20, sin81_30, sin81_40, sin81_50},
	{sin82_00, sin82_10, sin82_20, sin82_30, sin82_40, sin82_50},
	{sin83_00, sin83_10, sin83_20, sin83_30, sin83_40, sin83_50},
	{sin84_00, sin84_10, sin84_20, sin84_30, sin84_40, sin84_50},
	{sin85_00, sin85_10, sin85_20, sin85_30, sin85_40, sin85_50},
	{sin86_00, sin86_10, sin86_20, sin86_30, sin86_40, sin86_50},
	{sin87_00, sin87_10, sin87_20, sin87_30, sin87_40, sin87_50},
	{sin88_00, sin88_10, sin88_20, sin88_30, sin88_40, sin88_50},
	{sin89_00, sin89_10, sin89_20, sin89_30, sin89_40, sin89_50},
	{sin90_00, -1, -1, -1, -1},
}

var (
	tan_0_00  int32 = 0
	tan_0_10  int32 = 29
	tan_0_20  int32 = 58
	tan_0_30  int32 = 87
	tan_0_40  int32 = 116
	tan_0_50  int32 = 145
	tan_1_00  int32 = 175
	tan_1_10  int32 = 204
	tan_1_20  int32 = 233
	tan_1_30  int32 = 262
	tan_1_40  int32 = 291
	tan_1_50  int32 = 320
	tan_2_00  int32 = 349
	tan_2_10  int32 = 378
	tan_2_20  int32 = 407
	tan_2_30  int32 = 437
	tan_2_40  int32 = 466
	tan_2_50  int32 = 495
	tan_3_00  int32 = 524
	tan_3_10  int32 = 553
	tan_3_20  int32 = 582
	tan_3_30  int32 = 612
	tan_3_40  int32 = 641
	tan_3_50  int32 = 670
	tan_4_00  int32 = 699
	tan_4_10  int32 = 729
	tan_4_20  int32 = 758
	tan_4_30  int32 = 787
	tan_4_40  int32 = 816
	tan_4_50  int32 = 846
	tan_5_00  int32 = 875
	tan_5_10  int32 = 904
	tan_5_20  int32 = 934
	tan_5_30  int32 = 963
	tan_5_40  int32 = 992
	tan_5_50  int32 = 1022
	tan_6_00  int32 = 1051
	tan_6_10  int32 = 1080
	tan_6_20  int32 = 1110
	tan_6_30  int32 = 1139
	tan_6_40  int32 = 1169
	tan_6_50  int32 = 1198
	tan_7_00  int32 = 1228
	tan_7_10  int32 = 1257
	tan_7_20  int32 = 1287
	tan_7_30  int32 = 1317
	tan_7_40  int32 = 1346
	tan_7_50  int32 = 1376
	tan_8_00  int32 = 1405
	tan_8_10  int32 = 1435
	tan_8_20  int32 = 1465
	tan_8_30  int32 = 1495
	tan_8_40  int32 = 1524
	tan_8_50  int32 = 1554
	tan_9_00  int32 = 1584
	tan_9_10  int32 = 1614
	tan_9_20  int32 = 1644
	tan_9_30  int32 = 1673
	tan_9_40  int32 = 1703
	tan_9_50  int32 = 1733
	tan_10_00 int32 = 1763
	tan_10_10 int32 = 1793
	tan_10_20 int32 = 1823
	tan_10_30 int32 = 1853
	tan_10_40 int32 = 1883
	tan_10_50 int32 = 1914
	tan_11_00 int32 = 1944
	tan_11_10 int32 = 1974
	tan_11_20 int32 = 2004
	tan_11_30 int32 = 2035
	tan_11_40 int32 = 2065
	tan_11_50 int32 = 2095
	tan_12_00 int32 = 2126
	tan_12_10 int32 = 2156
	tan_12_20 int32 = 2186
	tan_12_30 int32 = 2217
	tan_12_40 int32 = 2247
	tan_12_50 int32 = 2278
	tan_13_00 int32 = 2309
	tan_13_10 int32 = 2339
	tan_13_20 int32 = 2370
	tan_13_30 int32 = 2401
	tan_13_40 int32 = 2432
	tan_13_50 int32 = 2462
	tan_14_00 int32 = 2493
	tan_14_10 int32 = 2524
	tan_14_20 int32 = 2555
	tan_14_30 int32 = 2586
	tan_14_40 int32 = 2617
	tan_14_50 int32 = 2648
	tan_15_00 int32 = 2679
	tan_15_10 int32 = 2711
	tan_15_20 int32 = 2742
	tan_15_30 int32 = 2773
	tan_15_40 int32 = 2805
	tan_15_50 int32 = 2836
	tan_16_00 int32 = 2867
	tan_16_10 int32 = 2899
	tan_16_20 int32 = 2931
	tan_16_30 int32 = 2962
	tan_16_40 int32 = 2994
	tan_16_50 int32 = 3026
	tan_17_00 int32 = 3057
	tan_17_10 int32 = 3089
	tan_17_20 int32 = 3121
	tan_17_30 int32 = 3153
	tan_17_40 int32 = 3185
	tan_17_50 int32 = 3217
	tan_18_00 int32 = 3249
	tan_18_10 int32 = 3281
	tan_18_20 int32 = 3314
	tan_18_30 int32 = 3346
	tan_18_40 int32 = 3378
	tan_18_50 int32 = 3411
	tan_19_00 int32 = 3443
	tan_19_10 int32 = 3476
	tan_19_20 int32 = 3508
	tan_19_30 int32 = 3541
	tan_19_40 int32 = 3574
	tan_19_50 int32 = 3607
	tan_20_00 int32 = 3640
	tan_20_10 int32 = 3673
	tan_20_20 int32 = 3706
	tan_20_30 int32 = 3739
	tan_20_40 int32 = 3772
	tan_20_50 int32 = 3805
	tan_21_00 int32 = 3839
	tan_21_10 int32 = 3872
	tan_21_20 int32 = 3906
	tan_21_30 int32 = 3939
	tan_21_40 int32 = 3973
	tan_21_50 int32 = 4006
	tan_22_00 int32 = 4040
	tan_22_10 int32 = 4074
	tan_22_20 int32 = 4108
	tan_22_30 int32 = 4142
	tan_22_40 int32 = 4176
	tan_22_50 int32 = 4210
	tan_23_00 int32 = 4245
	tan_23_10 int32 = 4279
	tan_23_20 int32 = 4314
	tan_23_30 int32 = 4348
	tan_23_40 int32 = 4383
	tan_23_50 int32 = 4417
	tan_24_00 int32 = 4452
	tan_24_10 int32 = 4487
	tan_24_20 int32 = 4522
	tan_24_30 int32 = 4557
	tan_24_40 int32 = 4592
	tan_24_50 int32 = 4628
	tan_25_00 int32 = 4663
	tan_25_10 int32 = 4699
	tan_25_20 int32 = 4734
	tan_25_30 int32 = 4770
	tan_25_40 int32 = 4806
	tan_25_50 int32 = 4841
	tan_26_00 int32 = 4877
	tan_26_10 int32 = 4913
	tan_26_20 int32 = 4950
	tan_26_30 int32 = 4986
	tan_26_40 int32 = 5022
	tan_26_50 int32 = 5059
	tan_27_00 int32 = 5095
	tan_27_10 int32 = 5132
	tan_27_20 int32 = 5169
	tan_27_30 int32 = 5206
	tan_27_40 int32 = 5243
	tan_27_50 int32 = 5280
	tan_28_00 int32 = 5317
	tan_28_10 int32 = 5354
	tan_28_20 int32 = 5392
	tan_28_30 int32 = 5430
	tan_28_40 int32 = 5467
	tan_28_50 int32 = 5505
	tan_29_00 int32 = 5543
	tan_29_10 int32 = 5581
	tan_29_20 int32 = 5619
	tan_29_30 int32 = 5658
	tan_29_40 int32 = 5696
	tan_29_50 int32 = 5735
	tan_30_00 int32 = 5774
	tan_30_10 int32 = 5812
	tan_30_20 int32 = 5851
	tan_30_30 int32 = 5890
	tan_30_40 int32 = 5930
	tan_30_50 int32 = 5969
	tan_31_00 int32 = 6009
	tan_31_10 int32 = 6048
	tan_31_20 int32 = 6088
	tan_31_30 int32 = 6128
	tan_31_40 int32 = 6168
	tan_31_50 int32 = 6208
	tan_32_00 int32 = 6249
	tan_32_10 int32 = 6289
	tan_32_20 int32 = 6330
	tan_32_30 int32 = 6371
	tan_32_40 int32 = 6412
	tan_32_50 int32 = 6453
	tan_33_00 int32 = 6494
	tan_33_10 int32 = 6536
	tan_33_20 int32 = 6577
	tan_33_30 int32 = 6619
	tan_33_40 int32 = 6661
	tan_33_50 int32 = 6703
	tan_34_00 int32 = 6745
	tan_34_10 int32 = 6787
	tan_34_20 int32 = 6830
	tan_34_30 int32 = 6873
	tan_34_40 int32 = 6916
	tan_34_50 int32 = 6959
	tan_35_00 int32 = 7002
	tan_35_10 int32 = 7046
	tan_35_20 int32 = 7089
	tan_35_30 int32 = 7133
	tan_35_40 int32 = 7177
	tan_35_50 int32 = 7221
	tan_36_00 int32 = 7265
	tan_36_10 int32 = 7310
	tan_36_20 int32 = 7355
	tan_36_30 int32 = 7400
	tan_36_40 int32 = 7445
	tan_36_50 int32 = 7490
	tan_37_00 int32 = 7536
	tan_37_10 int32 = 7581
	tan_37_20 int32 = 7627
	tan_37_30 int32 = 7673
	tan_37_40 int32 = 7720
	tan_37_50 int32 = 7766
	tan_38_00 int32 = 7813
	tan_38_10 int32 = 7860
	tan_38_20 int32 = 7907
	tan_38_30 int32 = 7954
	tan_38_40 int32 = 8002
	tan_38_50 int32 = 8050
	tan_39_00 int32 = 8098
	tan_39_10 int32 = 8146
	tan_39_20 int32 = 8195
	tan_39_30 int32 = 8243
	tan_39_40 int32 = 8292
	tan_39_50 int32 = 8342
	tan_40_00 int32 = 8391
	tan_40_10 int32 = 8441
	tan_40_20 int32 = 8491
	tan_40_30 int32 = 8541
	tan_40_40 int32 = 8591
	tan_40_50 int32 = 8642
	tan_41_00 int32 = 8693
	tan_41_10 int32 = 8744
	tan_41_20 int32 = 8796
	tan_41_30 int32 = 8847
	tan_41_40 int32 = 8899
	tan_41_50 int32 = 8952
	tan_42_00 int32 = 9004
	tan_42_10 int32 = 9057
	tan_42_20 int32 = 9110
	tan_42_30 int32 = 9163
	tan_42_40 int32 = 9217
	tan_42_50 int32 = 9271
	tan_43_00 int32 = 9325
	tan_43_10 int32 = 9380
	tan_43_20 int32 = 9435
	tan_43_30 int32 = 9490
	tan_43_40 int32 = 9545
	tan_43_50 int32 = 9601
	tan_44_00 int32 = 9657
	tan_44_10 int32 = 9713
	tan_44_20 int32 = 9770
	tan_44_30 int32 = 9827
	tan_44_40 int32 = 9884
	tan_44_50 int32 = 9942
	tan_45_00 int32 = 10000
	tan_45_10 int32 = 10058
	tan_45_20 int32 = 10117
	tan_45_30 int32 = 10176
	tan_45_40 int32 = 10235
	tan_45_50 int32 = 10295
	tan_46_00 int32 = 10355
	tan_46_10 int32 = 10416
	tan_46_20 int32 = 10477
	tan_46_30 int32 = 10538
	tan_46_40 int32 = 10599
	tan_46_50 int32 = 10661
	tan_47_00 int32 = 10724
	tan_47_10 int32 = 10786
	tan_47_20 int32 = 10850
	tan_47_30 int32 = 10913
	tan_47_40 int32 = 10977
	tan_47_50 int32 = 11041
	tan_48_00 int32 = 11106
	tan_48_10 int32 = 11171
	tan_48_20 int32 = 11237
	tan_48_30 int32 = 11303
	tan_48_40 int32 = 11369
	tan_48_50 int32 = 11436
	tan_49_00 int32 = 11504
	tan_49_10 int32 = 11571
	tan_49_20 int32 = 11640
	tan_49_30 int32 = 11708
	tan_49_40 int32 = 11778
	tan_49_50 int32 = 11847
	tan_50_00 int32 = 11918
	tan_50_10 int32 = 11988
	tan_50_20 int32 = 12059
	tan_50_30 int32 = 12131
	tan_50_40 int32 = 12203
	tan_50_50 int32 = 12276
	tan_51_00 int32 = 12349
	tan_51_10 int32 = 12423
	tan_51_20 int32 = 12497
	tan_51_30 int32 = 12572
	tan_51_40 int32 = 12647
	tan_51_50 int32 = 12723
	tan_52_00 int32 = 12799
	tan_52_10 int32 = 12876
	tan_52_20 int32 = 12954
	tan_52_30 int32 = 13032
	tan_52_40 int32 = 13111
	tan_52_50 int32 = 13190
	tan_53_00 int32 = 13270
	tan_53_10 int32 = 13351
	tan_53_20 int32 = 13432
	tan_53_30 int32 = 13514
	tan_53_40 int32 = 13597
	tan_53_50 int32 = 13680
	tan_54_00 int32 = 13764
	tan_54_10 int32 = 13848
	tan_54_20 int32 = 13934
	tan_54_30 int32 = 14019
	tan_54_40 int32 = 14106
	tan_54_50 int32 = 14193
	tan_55_00 int32 = 14281
	tan_55_10 int32 = 14370
	tan_55_20 int32 = 14460
	tan_55_30 int32 = 14550
	tan_55_40 int32 = 14641
	tan_55_50 int32 = 14733
	tan_56_00 int32 = 14826
	tan_56_10 int32 = 14919
	tan_56_20 int32 = 15013
	tan_56_30 int32 = 15108
	tan_56_40 int32 = 15204
	tan_56_50 int32 = 15301
	tan_57_00 int32 = 15399
	tan_57_10 int32 = 15497
	tan_57_20 int32 = 15597
	tan_57_30 int32 = 15697
	tan_57_40 int32 = 15798
	tan_57_50 int32 = 15900
	tan_58_00 int32 = 16003
	tan_58_10 int32 = 16107
	tan_58_20 int32 = 16212
	tan_58_30 int32 = 16319
	tan_58_40 int32 = 16426
	tan_58_50 int32 = 16534
	tan_59_00 int32 = 16643
	tan_59_10 int32 = 16753
	tan_59_20 int32 = 16864
	tan_59_30 int32 = 16977
	tan_59_40 int32 = 17090
	tan_59_50 int32 = 17205
	tan_60_00 int32 = 17321
	tan_60_10 int32 = 17437
	tan_60_20 int32 = 17556
	tan_60_30 int32 = 17675
	tan_60_40 int32 = 17796
	tan_60_50 int32 = 17917
	tan_61_00 int32 = 18040
	tan_61_10 int32 = 18165
	tan_61_20 int32 = 18291
	tan_61_30 int32 = 18418
	tan_61_40 int32 = 18546
	tan_61_50 int32 = 18676
	tan_62_00 int32 = 18807
	tan_62_10 int32 = 18940
	tan_62_20 int32 = 19074
	tan_62_30 int32 = 19210
	tan_62_40 int32 = 19347
	tan_62_50 int32 = 19486
	tan_63_00 int32 = 19626
	tan_63_10 int32 = 19768
	tan_63_20 int32 = 19912
	tan_63_30 int32 = 20057
	tan_63_40 int32 = 20204
	tan_63_50 int32 = 20353
	tan_64_00 int32 = 20503
	tan_64_10 int32 = 20655
	tan_64_20 int32 = 20809
	tan_64_30 int32 = 20965
	tan_64_40 int32 = 21123
	tan_64_50 int32 = 21283
	tan_65_00 int32 = 21445
	tan_65_10 int32 = 21609
	tan_65_20 int32 = 21775
	tan_65_30 int32 = 21943
	tan_65_40 int32 = 22113
	tan_65_50 int32 = 22286
	tan_66_00 int32 = 22460
	tan_66_10 int32 = 22637
	tan_66_20 int32 = 22817
	tan_66_30 int32 = 22998
	tan_66_40 int32 = 23183
	tan_66_50 int32 = 23369
	tan_67_00 int32 = 23559
	tan_67_10 int32 = 23750
	tan_67_20 int32 = 23945
	tan_67_30 int32 = 24142
	tan_67_40 int32 = 24342
	tan_67_50 int32 = 24545
	tan_68_00 int32 = 24751
	tan_68_10 int32 = 24960
	tan_68_20 int32 = 25172
	tan_68_30 int32 = 25386
	tan_68_40 int32 = 25605
	tan_68_50 int32 = 25826
	tan_69_00 int32 = 26051
	tan_69_10 int32 = 26279
	tan_69_20 int32 = 26511
	tan_69_30 int32 = 26746
	tan_69_40 int32 = 26985
	tan_69_50 int32 = 27228
	tan_70_00 int32 = 27475
	tan_70_10 int32 = 27725
	tan_70_20 int32 = 27980
	tan_70_30 int32 = 28239
	tan_70_40 int32 = 28502
	tan_70_50 int32 = 28770
	tan_71_00 int32 = 29042
	tan_71_10 int32 = 29319
	tan_71_20 int32 = 29600
	tan_71_30 int32 = 29887
	tan_71_40 int32 = 30178
	tan_71_50 int32 = 30475
	tan_72_00 int32 = 30777
	tan_72_10 int32 = 31084
	tan_72_20 int32 = 31397
	tan_72_30 int32 = 31761
	tan_72_40 int32 = 32041
	tan_72_50 int32 = 32371
	tan_73_00 int32 = 32709
	tan_73_10 int32 = 33052
	tan_73_20 int32 = 33402
	tan_73_30 int32 = 33759
	tan_73_40 int32 = 34124
	tan_73_50 int32 = 34495
	tan_74_00 int32 = 34874
	tan_74_10 int32 = 35261
	tan_74_20 int32 = 35656
	tan_74_30 int32 = 36059
	tan_74_40 int32 = 36470
	tan_74_50 int32 = 36891
	tan_75_00 int32 = 37321
	tan_75_10 int32 = 37760
	tan_75_20 int32 = 38208
	tan_75_30 int32 = 38667
	tan_75_40 int32 = 39136
	tan_75_50 int32 = 39617
	tan_76_00 int32 = 40108
	tan_76_10 int32 = 40611
	tan_76_20 int32 = 41126
	tan_76_30 int32 = 41653
	tan_76_40 int32 = 42193
	tan_76_50 int32 = 42747
	tan_77_00 int32 = 43315
	tan_77_10 int32 = 43897
	tan_77_20 int32 = 44494
	tan_77_30 int32 = 45107
	tan_77_40 int32 = 45736
	tan_77_50 int32 = 46382
	tan_78_00 int32 = 47046
	tan_78_10 int32 = 47729
	tan_78_20 int32 = 48430
	tan_78_30 int32 = 48152
	tan_78_40 int32 = 49894
	tan_78_50 int32 = 50658
	tan_79_00 int32 = 51446
	tan_79_10 int32 = 52257
	tan_79_20 int32 = 53093
	tan_79_30 int32 = 53955
	tan_79_40 int32 = 54845
	tan_79_50 int32 = 55764
	tan_80_00 int32 = 56713
	tan_80_10 int32 = 57694
	tan_80_20 int32 = 58708
	tan_80_30 int32 = 59758
	tan_80_40 int32 = 60844
	tan_80_50 int32 = 61970
	tan_81_00 int32 = 63138
	tan_81_10 int32 = 64348
	tan_81_20 int32 = 65606
	tan_81_30 int32 = 66912
	tan_81_40 int32 = 68269
	tan_81_50 int32 = 69682
	tan_82_00 int32 = 71154
	tan_82_10 int32 = 72687
	tan_82_20 int32 = 74287
	tan_82_30 int32 = 75958
	tan_82_40 int32 = 77704
	tan_82_50 int32 = 79530
	tan_83_00 int32 = 81443
	tan_83_10 int32 = 83450
	tan_83_20 int32 = 85555
	tan_83_30 int32 = 87769
	tan_83_40 int32 = 90098
	tan_83_50 int32 = 92553
	tan_84_00 int32 = 95144
	tan_84_10 int32 = 97882
	tan_84_20 int32 = 100780
	tan_84_30 int32 = 103854
	tan_84_40 int32 = 107119
	tan_84_50 int32 = 110594
	tan_85_00 int32 = 114301
	tan_85_10 int32 = 118262
	tan_85_20 int32 = 122505
	tan_85_30 int32 = 127062
	tan_85_40 int32 = 131969
	tan_85_50 int32 = 137267
	tan_86_00 int32 = 143007
	tan_86_10 int32 = 149244
	tan_86_20 int32 = 156048
	tan_86_30 int32 = 163499
	tan_86_40 int32 = 171693
	tan_86_50 int32 = 180750
	tan_87_00 int32 = 190811
	tan_87_10 int32 = 202056
	tan_87_20 int32 = 214704
	tan_87_30 int32 = 229038
	tan_87_40 int32 = 245418
	tan_87_50 int32 = 264316
	tan_88_00 int32 = 286363
	tan_88_10 int32 = 312416
	tan_88_20 int32 = 343678
	tan_88_30 int32 = 381885
	tan_88_40 int32 = 429641
	tan_88_50 int32 = 491039
	tan_89_00 int32 = 572900
	tan_89_10 int32 = 687501
	tan_89_20 int32 = 859398
	tan_89_30 int32 = 1145887
	tan_89_40 int32 = 1718854
	tan_89_50 int32 = 3437737
	tan_90_00 int32 = math.MaxInt32
)

var tanval = [][6]int32{
	{tan_0_00, tan_0_10, tan_0_20, tan_0_30, tan_0_40, tan_0_50},
	{tan_1_00, tan_1_10, tan_1_20, tan_1_30, tan_1_40, tan_1_50},
	{tan_2_00, tan_2_10, tan_2_20, tan_2_30, tan_2_40, tan_2_50},
	{tan_3_00, tan_3_10, tan_3_20, tan_3_30, tan_3_40, tan_3_50},
	{tan_4_00, tan_4_10, tan_4_20, tan_4_30, tan_4_40, tan_4_50},
	{tan_5_00, tan_5_10, tan_5_20, tan_5_30, tan_5_40, tan_5_50},
	{tan_6_00, tan_6_10, tan_6_20, tan_6_30, tan_6_40, tan_6_50},
	{tan_7_00, tan_7_10, tan_7_20, tan_7_30, tan_7_40, tan_7_50},
	{tan_8_00, tan_8_10, tan_8_20, tan_8_30, tan_8_40, tan_8_50},
	{tan_9_00, tan_9_10, tan_9_20, tan_9_30, tan_9_40, tan_9_50},
	{tan_10_00, tan_10_10, tan_10_20, tan_10_30, tan_10_40, tan_10_50}, // 10
	{tan_11_00, tan_11_10, tan_11_20, tan_11_30, tan_11_40, tan_11_50},
	{tan_12_00, tan_12_10, tan_12_20, tan_12_30, tan_12_40, tan_12_50},
	{tan_13_00, tan_13_10, tan_13_20, tan_13_30, tan_13_40, tan_13_50},
	{tan_14_00, tan_14_10, tan_14_20, tan_14_30, tan_14_40, tan_14_50},
	{tan_15_00, tan_15_10, tan_15_20, tan_15_30, tan_15_40, tan_15_50},
	{tan_16_00, tan_16_10, tan_16_20, tan_16_30, tan_16_40, tan_16_50},
	{tan_17_00, tan_17_10, tan_17_20, tan_17_30, tan_17_40, tan_17_50},
	{tan_18_00, tan_18_10, tan_18_20, tan_18_30, tan_18_40, tan_18_50},
	{tan_19_00, tan_19_10, tan_19_20, tan_19_30, tan_19_40, tan_19_50},
	{tan_20_00, tan_20_10, tan_20_20, tan_20_30, tan_20_40, tan_20_50}, // 20
	{tan_21_00, tan_21_10, tan_21_20, tan_21_30, tan_21_40, tan_21_50},
	{tan_22_00, tan_22_10, tan_22_20, tan_22_30, tan_22_40, tan_22_50},
	{tan_23_00, tan_23_10, tan_23_20, tan_23_30, tan_23_40, tan_23_50},
	{tan_24_00, tan_24_10, tan_24_20, tan_24_30, tan_24_40, tan_24_50},
	{tan_25_00, tan_25_10, tan_25_20, tan_25_30, tan_25_40, tan_25_50},
	{tan_26_00, tan_26_10, tan_26_20, tan_26_30, tan_26_40, tan_26_50},
	{tan_27_00, tan_27_10, tan_27_20, tan_27_30, tan_27_40, tan_27_50},
	{tan_28_00, tan_28_10, tan_28_20, tan_28_30, tan_28_40, tan_28_50},
	{tan_29_00, tan_29_10, tan_29_20, tan_29_30, tan_29_40, tan_29_50},
	{tan_30_00, tan_30_10, tan_30_20, tan_30_30, tan_30_40, tan_30_50}, // 30
	{tan_31_00, tan_31_10, tan_31_20, tan_31_30, tan_31_40, tan_31_50},
	{tan_32_00, tan_32_10, tan_32_20, tan_32_30, tan_32_40, tan_32_50},
	{tan_33_00, tan_33_10, tan_33_20, tan_33_30, tan_33_40, tan_33_50},
	{tan_34_00, tan_34_10, tan_34_20, tan_34_30, tan_34_40, tan_34_50},
	{tan_35_00, tan_35_10, tan_35_20, tan_35_30, tan_35_40, tan_35_50},
	{tan_36_00, tan_36_10, tan_36_20, tan_36_30, tan_36_40, tan_36_50},
	{tan_37_00, tan_37_10, tan_37_20, tan_37_30, tan_37_40, tan_37_50},
	{tan_38_00, tan_38_10, tan_38_20, tan_38_30, tan_38_40, tan_38_50},
	{tan_39_00, tan_39_10, tan_39_20, tan_39_30, tan_39_40, tan_39_50},
	{tan_40_00, tan_40_10, tan_40_20, tan_40_30, tan_40_40, tan_40_50}, // 40
	{tan_41_00, tan_41_10, tan_41_20, tan_41_30, tan_41_40, tan_41_50},
	{tan_42_00, tan_42_10, tan_42_20, tan_42_30, tan_42_40, tan_42_50},
	{tan_43_00, tan_43_10, tan_43_20, tan_43_30, tan_43_40, tan_43_50},
	{tan_44_00, tan_44_10, tan_44_20, tan_44_30, tan_44_40, tan_44_50},
	{tan_45_00, tan_45_10, tan_45_20, tan_45_30, tan_45_40, tan_45_50},
	{tan_46_00, tan_46_10, tan_46_20, tan_46_30, tan_46_40, tan_46_50},
	{tan_47_00, tan_47_10, tan_47_20, tan_47_30, tan_47_40, tan_47_50},
	{tan_48_00, tan_48_10, tan_48_20, tan_48_30, tan_48_40, tan_48_50},
	{tan_49_00, tan_49_10, tan_49_20, tan_49_30, tan_49_40, tan_49_50},
	{tan_50_00, tan_50_10, tan_50_20, tan_50_30, tan_50_40, tan_50_50}, // 50
	{tan_51_00, tan_51_10, tan_51_20, tan_51_30, tan_51_40, tan_51_50},
	{tan_52_00, tan_52_10, tan_52_20, tan_52_30, tan_52_40, tan_52_50},
	{tan_53_00, tan_53_10, tan_53_20, tan_53_30, tan_53_40, tan_53_50},
	{tan_54_00, tan_54_10, tan_54_20, tan_54_30, tan_54_40, tan_54_50},
	{tan_55_00, tan_55_10, tan_55_20, tan_55_30, tan_55_40, tan_55_50},
	{tan_56_00, tan_56_10, tan_56_20, tan_56_30, tan_56_40, tan_56_50},
	{tan_57_00, tan_57_10, tan_57_20, tan_57_30, tan_57_40, tan_57_50},
	{tan_58_00, tan_58_10, tan_58_20, tan_58_30, tan_58_40, tan_58_50},
	{tan_59_00, tan_59_10, tan_59_20, tan_59_30, tan_59_40, tan_59_50},
	{tan_60_00, tan_60_10, tan_60_20, tan_60_30, tan_60_40, tan_60_50}, // 60
	{tan_61_00, tan_61_10, tan_61_20, tan_61_30, tan_61_40, tan_61_50},
	{tan_62_00, tan_62_10, tan_62_20, tan_62_30, tan_62_40, tan_62_50},
	{tan_63_00, tan_63_10, tan_63_20, tan_63_30, tan_63_40, tan_63_50},
	{tan_64_00, tan_64_10, tan_64_20, tan_64_30, tan_64_40, tan_64_50},
	{tan_65_00, tan_65_10, tan_65_20, tan_65_30, tan_65_40, tan_65_50},
	{tan_66_00, tan_66_10, tan_66_20, tan_66_30, tan_66_40, tan_66_50},
	{tan_67_00, tan_67_10, tan_67_20, tan_67_30, tan_67_40, tan_67_50},
	{tan_68_00, tan_68_10, tan_68_20, tan_68_30, tan_68_40, tan_68_50},
	{tan_69_00, tan_69_10, tan_69_20, tan_69_30, tan_69_40, tan_69_50},
	{tan_70_00, tan_70_10, tan_70_20, tan_70_30, tan_70_40, tan_70_50}, // 70
	{tan_71_00, tan_71_10, tan_71_20, tan_71_30, tan_71_40, tan_71_50},
	{tan_72_00, tan_72_10, tan_72_20, tan_72_30, tan_72_40, tan_72_50},
	{tan_73_00, tan_73_10, tan_73_20, tan_73_30, tan_73_40, tan_73_50},
	{tan_74_00, tan_74_10, tan_74_20, tan_74_30, tan_74_40, tan_74_50},
	{tan_75_00, tan_75_10, tan_75_20, tan_75_30, tan_75_40, tan_75_50},
	{tan_76_00, tan_76_10, tan_76_20, tan_76_30, tan_76_40, tan_76_50},
	{tan_77_00, tan_77_10, tan_77_20, tan_77_30, tan_77_40, tan_77_50},
	{tan_78_00, tan_78_10, tan_78_20, tan_78_30, tan_78_40, tan_78_50},
	{tan_79_00, tan_79_10, tan_79_20, tan_79_30, tan_79_40, tan_79_50},
	{tan_80_00, tan_80_10, tan_80_20, tan_80_30, tan_80_40, tan_80_50}, // 80
	{tan_81_00, tan_81_10, tan_81_20, tan_81_30, tan_81_40, tan_81_50},
	{tan_82_00, tan_82_10, tan_82_20, tan_82_30, tan_82_40, tan_82_50},
	{tan_83_00, tan_83_10, tan_83_20, tan_83_30, tan_83_40, tan_83_50},
	{tan_84_00, tan_84_10, tan_84_20, tan_84_30, tan_84_40, tan_84_50},
	{tan_85_00, tan_85_10, tan_85_20, tan_85_30, tan_85_40, tan_85_50},
	{tan_86_00, tan_86_10, tan_86_20, tan_86_30, tan_86_40, tan_86_50},
	{tan_87_00, tan_87_10, tan_87_20, tan_87_30, tan_87_40, tan_87_50},
	{tan_88_00, tan_88_10, tan_88_20, tan_88_30, tan_88_40, tan_88_50},
	{tan_89_00, tan_89_10, tan_89_20, tan_89_30, tan_89_40, tan_89_50},
	{tan_90_00, tan_90_00, tan_90_00, tan_90_00, tan_90_00, tan_90_00}, // 90
}

// 正弦
func Sine(angle Angle) (int32, int32) {
	if angle.minute >= 60 {
		panic(fmt.Sprintf("base: invalid minute param %v for sin value", angle.minute))
	}
	if int32(angle.degree)*int32(angle.minute) < 0 {
		panic(fmt.Sprintf("base: invalid degree param %v and minute param %v for sin value", angle.degree, angle.minute))
	}
	angle.Normalize()
	if angle.degree >= 90 && angle.degree < 180 {
		angle.degree -= 90
		return Cosine(angle)
	} else if angle.degree >= 180 && angle.degree < 270 {
		angle.degree -= 180
		n, d := Sine(angle)
		return -n, d
	} else if angle.degree >= 270 && angle.degree < 360 {
		angle.degree -= 180
		n, d := Sine(angle)
		return -n, d
	}
	return sinval[angle.degree][angle.minute/10], denominator
}

// 餘弦
func Cosine(angle Angle) (int32, int32) {
	if angle.minute >= 60 {
		panic(fmt.Sprintf("base: invalid minute param %v for cos value", angle.minute))
	}
	if int32(angle.degree)*int32(angle.minute) < 0 {
		panic(fmt.Sprintf("base: invalid degree param %v and minute param %v", angle.degree, angle.minute))
	}
	angle.Normalize()
	if angle.degree >= 90 && angle.degree < 180 {
		angle.degree -= 90
		n, d := Sine(angle)
		return -n, d
	} else if angle.degree >= 180 && angle.degree < 270 {
		angle.degree -= 180
		n, d := Cosine(angle)
		return -n, d
	} else if angle.degree >= 270 && angle.degree < 360 {
		angle.degree -= 180
		n, d := Cosine(angle)
		return -n, d
	}

	dm := 90*60 - (angle.degree*60 + angle.minute)
	angle.degree, angle.minute = dm/60, dm%60
	return sinval[angle.degree][angle.minute/10], denominator
}

// 正切
func Tangent(angle Angle) (int32, int32) {
	if angle.minute >= 60 {
		panic(fmt.Sprintf("base: invalid minute param %v for tangent value", angle.minute))
	}
	if int32(angle.degree)*int32(angle.minute) < 0 {
		panic(fmt.Sprintf("base: invalid degree param %v and minute param %v for tangent angle", angle.degree, angle.minute))
	}

	var negative bool
	angle.Normalize()
	if angle.degree >= 90 && angle.degree < 180 { // 第二象限 x<=0  y>=0
		angle.degree -= 90
		n, d := Cotangent(angle)
		return n, -d //
	} else if angle.degree >= 180 && angle.degree < 270 { // 第三象限  x<=0 y<=0
		angle.degree -= 180
		negative = true
	} else if angle.degree >= 270 && angle.degree < 360 { // 第四象限  x>=0 y<=0
		tp := TwoPiAngle()
		tp.Sub(angle)
		n, d := Tangent(tp)
		return -n, d
	}
	if negative {
		return -tanval[angle.degree][angle.minute/10], -denominator
	}
	return tanval[angle.degree][angle.minute/10], denominator
}

// 餘切
func Cotangent(angle Angle) (int32, int32) {
	if angle.minute >= 60 {
		panic(fmt.Sprintf("base: invalid minute param %v for tangent value", angle.minute))
	}
	if int32(angle.degree)*int32(angle.minute) < 0 {
		panic(fmt.Sprintf("base: invalid degree param %v and minute param %v for cotangent angle", angle.degree, angle.minute))
	}

	var negative bool
	angle.Normalize()
	if angle.degree >= 90 && angle.degree < 180 {
		angle.degree -= 90
		n, d := Tangent(angle)
		return -n, d
	} else if angle.degree >= 180 && angle.degree < 270 {
		angle.degree -= 180
		negative = true
	} else if angle.degree >= 270 && angle.degree < 360 {
		tp := TwoPiAngle()
		tp.Sub(angle)
		n, d := Cotangent(tp)
		return n, -d
	}

	if angle.degree == 90 && angle.minute == 0 {
		return 0, denominator
	}

	dm := 90*60 - angle.degree*60 - angle.minute
	angle.degree, angle.minute = dm/60, dm%60
	if negative {
		return -tanval[angle.degree][angle.minute], -denominator
	}
	return tanval[angle.degree][angle.minute/10], denominator
}

func ArcSine(sn, sd int32) Angle {
	if sd == 0 {
		panic(fmt.Sprintf("base: invalid denominator %v for ArcSine", sd))
	}
	if sn == 0 {
		return Angle{0, 0}
	}
	var (
		l, r     = int16(0), int16(len(sinval)) - 1
		m        = r >> 1
		n        int16
		negative bool
	)
	if sn < 0 && sd < 0 {
		sn = -sn
		sd = -sd
	} else if sn < 0 {
		sn = -sn
		negative = true
	} else if sd < 0 {
		sd = -sd
		negative = true
	}
	for l <= r {
		if sn*denominator < sd*sinval[m][0] {
			r = m - 1
			m = (l + r) >> 1
			n = 0
		} else if sn*denominator > sd*sinval[m][5] {
			l = m + 1
			m = (l + r) >> 1
			n = 5
		} else {
			goto bl
		}
	}
	if negative {
		m = -m
		n = -n
	}
	return Angle{m, n * 10}

bl:
	for i := int16(0); i < int16(len(sinval[m])); i++ {
		if sn*denominator <= sd*sinval[m][i] {
			n = i
			break
		}
	}
	if negative {
		m = -m
		n = -n
	}
	return Angle{m, n * 10}
}

// 反餘弦
func ArcCosine(cn, cd int32) Angle {
	angle := ArcSine(cn, cd)
	minutes := 90*60 - (angle.degree*60 + angle.minute)
	return Angle{degree: minutes / 60, minute: minutes % 60}
}

// 反正切
func ArcTangent(y, x int32) Angle {
	if x == 0 {
		if y < 0 {
			return Angle{270, 0}
		} else if y > 0 {
			return Angle{90, 0}
		}
	}
	if y == 0 {
		if x < 0 {
			return Angle{180, 0}
		} else if x > 0 {
			return Angle{0, 0}
		}
	}

	var (
		l, r     = int16(0), int16(len(tanval)) - 1
		m        = r >> 1
		n        int16
		negative bool
	)
	if x < 0 && y < 0 {
		x = -x
		y = -y
	} else if x < 0 {
		x = -x
		negative = true
	} else if y < 0 {
		y = -y
		negative = true
	}
	for l <= r {
		if y*denominator < x*tanval[m][0] {
			r = m - 1
			m = (l + r) >> 1
			n = 0
		} else if y*denominator > x*tanval[m][5] {
			l = m + 1
			m = (l + r) >> 1
			n = 5
		} else {
			goto bl
		}
	}
	if negative {
		m = -m
		n = -n
	}
	return Angle{m, n * 10}

bl:
	for i := int16(0); i < int16(len(tanval[m])); i++ {
		if y*denominator <= x*tanval[m][i] {
			n = i
			break
		}
	}
	if negative {
		m = -m
		n = -n
	}
	return Angle{m, n * 10}
}

// 反餘切
func ArcCotangent(y, x int32) Angle {
	angle := ArcTangent(y, x)
	minutes := 90*60 - (angle.degree*60 + angle.minute)
	return Angle{degree: minutes / 60, minute: minutes % 60}
}

// 平方根計算
func Sqrt(num uint32) uint32 {
	if num <= 1 {
		return num
	}
	s := 1
	num1 := num - 1
	if num1 > 65535 {
		s += 8
		num1 >>= 16
	}
	if num1 > 255 {
		s += 4
		num1 >>= 8
	}
	if num1 > 15 {
		s += 2
		num1 >>= 4
	}
	if num1 > 3 {
		s += 1
	}
	x0 := uint32(1 << s)
	x1 := (x0 + (num >> s)) >> 1
	for x1 < x0 {
		x0 = x1
		x1 = (x0 + (num / x0)) >> 1
	}
	return x0
}
