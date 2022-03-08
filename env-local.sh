#!/bin/bash

# export $(xargs <.env.local)

export $(cat .env.local | xargs)