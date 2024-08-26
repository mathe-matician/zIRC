#!/bin/bash

cd $HOME
pg_ctl -D ./postgres/ -l logfile start