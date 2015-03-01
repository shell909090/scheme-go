#!/usr/bin/python
# -*- coding: utf-8 -*-
'''
@date: 2013-06-30
@author: shell.xu
'''
import os, sys

builtin={}
def define(name, evaled=None):
    def inner(func):
        if evaled is not None: func.evaled = evaled
        builtin[name] = func
        returnfunc
    return inner

@define(u'+')
def add(stack, envs, objs):
    pass
