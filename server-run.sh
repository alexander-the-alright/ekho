# ==============================================================================
# Auth: Alex Celani
# File: server-run.sh
# Revn: 10-05-2023  1.0
# Func: keep ekho server running at all times
#
# TODO: fucking figure this shit out, chief
# ==============================================================================
# CHANGE LOG
# ------------------------------------------------------------------------------
# 10-05-2023: thought
#
# ==============================================================================

# if I do this
# ~/bin/ekho-server &
# I have no way I of restarting the process when it ends
# might need to look into using Signals to restart, may PID?

# Other option is 
# for   # pseudo code, syntax isn't legit
#   ~/bin/ekho-server
#   ./server-run.sh
# self referential
# problem being, doing this will tie up the user, must be done on
# login by a kaz-like auto-login user that doesn't do anything other
# than run some helpful processes

echo "lmfao dude"
