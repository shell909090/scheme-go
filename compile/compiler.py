#!/usr/bin/python
# -*- coding: utf-8 -*-
'''
@date: 2013-06-30
@author: shell.xu
'''
import os, sys

class Frame(object):

    def add_code(self):
        pass

    def eval(self, obj):
        if isinstance(obj, OCons):
            
        

def compiler(objs):
    global_frame = Frame()
    for o in objs:
        global_frame.eval(o)
