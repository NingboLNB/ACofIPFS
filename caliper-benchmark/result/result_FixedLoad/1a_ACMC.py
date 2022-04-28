#!/usr/bin/env python
# coding=utf-8
import numpy as np
import matplotlib.pyplot as plt

from matplotlib.font_manager import FontProperties
font1 = {'family' : 'Times New Roman',
'weight' : 'normal',
'size' : 13,
}
font2 = {'family' : 'Times New Roman',
'weight' : 'normal',
'size' : 18,
}
def getChineseFont():
    return FontProperties(fname='/System/Library/Fonts/Times.ttc')


nodes = [100,150,200,250,300,350,400,450,500,550,600]
queryPolicy = [11.9,61.2,70.4,76.3,83.2,85.1,87.1,86.9,88.1,88.3,86.8]
queryObject = [11.6,63.7,70.3,78.1,86.6,86.7,86.4,88.7,88.5,89.4,90.3]
AddResource = [12.5,57.0,70.5,78.0,82.2,82.3,81.6,83.6,83.1,82.5,82.9]
AddPolicy = [11.3,55.2,71.3,78.4,80.3,79.9,80.1,80.2,80.7,80.5,80.2]
getToken = [12.1,52.4,62.1,68.5,72.8,73.6,72.2,76.1,73.5,75.7,73.5]
#TokenValidation = [11.6,57.0,66.9,76.9,80.4,80.1,78.5,83.8,83.1,83.7,80.5]
A,=plt.plot(nodes,queryObject,label="QueryObject",linewidth=1.5,color='dodgerblue',marker='s',ls='--',ms=4)
B,=plt.plot(nodes,queryPolicy,label="QueryPolicy",linewidth=1.5,color='tomato',marker='v',ls='-.',ms=4)
C,=plt.plot(nodes,AddResource,label="AddResource",linewidth=1.5,color='mediumslateblue',marker='x',ls='-.',ms=4)
D,=plt.plot(nodes,AddPolicy,label="AddPolicy",linewidth=1.5,color='darksalmon',marker='*',ls='--',ms=4)
E,=plt.plot(nodes,getToken,label="GetToken",linewidth=1.5,color='limegreen',marker='o',ls='--',ms=4)
#F,=plt.plot(nodes,TokenValidation,label="TokenValidation",linewidth=1.5,color='g',marker='o',ls='--',ms=4)




legend=plt.legend(handles=[A,B,C,D,E],prop=font1)


#plt.xlabel('Running Time/s\n(c) Convergence when the 80 task nodes',fontproperties=getChineseFont())
#plt.xlabel('The number of tasks\n',fontproperties=getChineseFont())
#plt.ylabel('Average Speedup',FontProperties=getChineseFont())
plt.xlabel('The Value of Fixed Load\n',font2)
plt.ylabel('Throughput',font2)

#plt.legend()
plt.savefig('/Users/liuningbo/Desktop/elsarticle-template_2/figure_6b.png')
plt.show()
