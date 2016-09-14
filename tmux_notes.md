# Tmux Conf File

You can have this in your // ~/.tmux.conf // file or have // /etc/tmux.conf // file for global configurations

	
	# change default action key
	set -g prefix C-a
	
	bind-key C-a last-window
	
	# change split window bindings
	unbind-key %
	unbind-key '"'
	bind-key -r | split-window -h
	bind-key -r - split-window -v
	bind-key -r > resize-pane -R 4
	bind-key -r < resize-pane -L 4
	bind-key -r ^ resize-pane -U 4
	bind-key -r v resize-pane -D 4
	
	#set -g mouse-select-pane on
	# automatically set window title
	setw -g automatic-rename


Basic premise is once you enter "command mode" (Ctrl+a) you can type to set the following

**|** Splits window Horizontally
**-** Vertical Split
**>** Resize to right
**<** Resize to left
**^** Resize upwards
**v** Resize downwards

# Tmux Commands

Open a New session just by typing "tmux"

	
	root@host# tmux


You can use your "window manipulation" commands you set up previously.

__General Commands__

 | Command  | Notes                          | 
 | -------  | -----                          | 
 | Ctrl+a   | Enters you into "command mode" | 
 | Ctrl+a c | Creates New window Session     | 
 | Ctrl+a n | Switch between window session  | 
 | Ctrl+a d | Disconnects from session       | 

Once you have disconnected you can list your actively running sessions.

	
	root@host# tmux ls
	0: 1 windows (created Tue Jan 22 10:02:25 2013) [172x38]
	1: 2 windows (created Tue Jan 22 10:02:39 2013) [172x38]


You can connect to individual sessions

	
	root@host# tmux attach -t 1


You can terminate your tmux sessions just by connecting to them; then using the standard "logout" Unix command.

If you get the following error

	
	root@host# tmux ls
	failed to connect to server: Connection refused


It just means there is no active tmux sessions

__Notes About 'Read Only'__\\
In the first terminal, start tmux where shared is the session name and shareds is the name of the socket (you can also use the shortcut *new* instead of *new-session*):

` tmux -S /tmp/shareds new-session -s shared `

Then chgrp the socket to a group that both users share in common. In this example, joint is the group that both users share. If there are other users in the group, then they also have access. So it might be recommended that the group have only the two members.

 ` chgrp joint /tmp/shareds `

In the second terminal attach using that socket and session.

` tmux -S /tmp/shareds attach -t shared `

That's it. The session can be made read-only for the second user, but only on a voluntary basis. The decision to work read-only is made when the second user attaches to the session.

` tmux -S /tmp/shareds attach -t shared -r `

__Parallel Commands__\\
If you have a Tmux window divided into panes, you can use the // synchronize-panes // window option to send each pane the same keyboard input simultaneously. You can do this by switching to the appropriate window, typing your Tmux prefix (Ctrl-A) and then a colon to bring up a Tmux command line, and type...

	
	:setw synchronize-panes on

