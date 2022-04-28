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
queryAR = [11.6,65.7,71.2,77.1,85.6,88.7,86.4,88.7,89.5,89.4,91.2]
queryVR = [11.9,64.2,70.4,76.5,86.2,88.1,87.1,86.9,90.1,92.1,91.7]
tokenValidation = [11.6,57.0,66.9,76.9,80.4,80.1,79.5,83.8,83.1,83.7,82.5]
B,=plt.plot(nodes,queryAR,label="GetAuthorizationLog",linewidth=1.5,color='limegreen',marker='v',ls='-.',ms=4)
A,=plt.plot(nodes,queryVR,label="GetVisitLog",linewidth=1.5,color='tomato',marker='s',ls='--',ms=4)
C,=plt.plot(nodes,tokenValidation,label="VerifyToken",linewidth=1.5,color='dodgerblue',marker='x',ls='-.',ms=4)
#D,=plt.plot(nodes,AddPolicy,label="AddPolicy",linewidth=1.5,color='b',marker='*',ls='--',ms=4)
#E,=plt.plot(nodes,GetToken,label="GetToken",linewidth=1.5,color='g',marker='o',ls='--',ms=4)




legend=plt.legend(handles=[B,A,C],prop=font1)


#plt.xlabel('Running Time/s\n(c) Convergence when the 80 task nodes',fontproperties=getChineseFont())
#plt.xlabel('The number of tasks\n',fontproperties=getChineseFont())
#plt.ylabel('Average Speedup',FontProperties=getChineseFont())
plt.xlabel('The Value of Fixed Load\n',font2)
plt.ylabel('Throughput',font2)

#plt.legend()
plt.savefig('/Users/liuningbo/Desktop/elsarticle-template_2/figure_6b.png')
plt.show()
