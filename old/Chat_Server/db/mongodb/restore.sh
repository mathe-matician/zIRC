#!/bin/bash
mongorestore -d ChatMessages /dump/ChatMessages
mongorestore -d ChatServerConfig /dump/ChatServerConfig
mongorestore -d ZIRC /dump/ZIRC