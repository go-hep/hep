#!/usr/bin/env python2

import numpy as np

gaussian = lambda x: 3*np.exp(-(30-x)**2/20.)

n = 100
xdata = np.arange(n)
ydata = gaussian(xdata)

import matplotlib.pyplot as plt

plt.plot(xdata, ydata, '.')

with open("gauss-data.txt","w") as f:
    for i in range(n):
        f.write("%+e %+e\n" % (xdata[i], ydata[i]))
        pass
    pass

X = np.arange(ydata.size)
x = np.sum(X*ydata)/np.sum(ydata)
width = np.sqrt(np.abs(np.sum((X-x)**2*ydata)/np.sum(ydata)))

ymax = ydata.max()

fit = lambda t : ymax*np.exp(-(t-x)**2/(2*width**2))

print "mean= %+e\nwidth= %+e" % (x, width)

plt.plot(fit(X), '-')
plt.show()
