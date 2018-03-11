# Git Notes

My git notes in no paticular order

* [Add your Keys](git_notes.md#add-your-keys)
* [Download Your Code](git_notes.md#download-your-code)
* [Committing](git_notes.md#committing)
* [Branching Out](git_notes.md#branching-out)
* [Aliases](git_notes.md#aliases)

## Add Your Keys

First and foremost you need a "git" and/or a "github" account ( http://github.com ). And you need to add your ssh key for the server you are going to use as your "workspace"

You can usually just copy/paste your `id_rsa.pub` contents.

```
root@host# cat ~/.ssh/id_rsa.pub
ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAIEA4f0SiWe7Evlu5S2NUfoEEX5gqUTInlB6Kwn7iWWAW8C7Zt2PFk9ZMGho2cUwV3cyNNxn4dKwwzv74tLTalmAstvCwfJZSYmXxDzKnbdHIH2kaWuyibMT9aHgYocRpfAf+TucRIB1yKyOHESk3XYREdprGGDG4MlhopTkIgEpy38= chrish@hostname
```

## Initialize Your Directory


You need to "initialize" a working directory where you will be working on your code. This is done simply by creating a working directory and using the `git init` command

```
root@host# mkdir /usr/local/git
root@host# cd /usr/local/git
root@host# git init
```

This creates various directories and files needed by git

## Download Your Code

Now you need to download your code into your working directory. This is done by "cloning" the code.

```
root@host# pwd
/usr/local/git
root@host# git clone git@git.4over.com:4over/sysops.git
```

**NOTE** It is always a good idea to make sure you have the LATEST codeset by syncing ("pulling") the most current copy

```
root@host# pwd
/usr/local/git
root@host# git pull
```

Now that you're here...you can `cd` into the directory that holds your code/project

```
root@host# pwd
/usr/local/git
root@host# ls -l
total 16K
drwxrwxr-x  3 chrish sysman 4.0K Dec 13 15:07 postgres/
drwxrwxr-x  4 chrish sysman 4.0K Dec 13 15:07 puppet/
-rw-rw-r--  1 chrish sysman   51 Dec 12 10:52 README.md
drwxrwxr-x 18 chrish sysman 4.0K Dec 12 11:13 ul4/
root@host# cd ul4
```

## Committing

Once you have finished editing your code/project; you "sync" it up with the repo in github.

In your working directory; "add" your changes (and yes you need the dot).

```
root@host# pwd
/usr/local/git
root@host# git add .
```

You can see what's pending commitment with the "status" command

```
root@host# git status
```

If you don't like what you see you can remove your pending commits with a "reset" command

```
root@host# git reset
```

When you are happy with what it's in your "pending commits" - you can commit the changes for syncing

```
root@host# git commit
```

This will open your favorite editor (VI of course ;-)) so you can add your comments
```
      1 This is a comment - do it on the first line
      2 # Please enter the commit message for your changes.
      3 # (Comment lines starting with '#' will not be included)
      4 # On branch master
      5 # Changes to be committed:
      6 #   (use "git reset HEAD <file>..." to unstage)
      7 #
      8 #       modified:   myfile.txt
      9 #
```

Now you can sync with the master repo by "pushing" your committed changes there.

```
root@host# git push
Counting objects: 7, done.
Compressing objects: 100% (4/4), done.
Writing objects: 100% (4/4), 383 bytes, done.
Total 4 (delta 3), reused 0 (delta 0)
To git@git.4over.com:4over/sysops.git
   35281da..546caf5  master -> master
```

## Branching Out

When using Git, it’s a good practice to do our work on a separate topic branch rather than the master branch

You do this with the following command (this "switches" you from working on the "master" to a new "copy" of the project named "static-page")

```
git checkout -b static-pages
```

Now all your `git add` and `git commit` commands will be on the `static-pages` branch. This allows you to work independently without "harming" the master.

When you are done and want to "merge" your changes into the master; you need to do the following

1) Make sure all your changes are committed (and/or pushed) to your branch

```
git add .
git commit -am "Finished static pages"
git push origin static-pages
```

2) "Switch" to the master branch by checking it out (you pass your branch name to the merge command)

```
git checkout master
git merge static-pages
```

3) Now push it into github

```
git push
```

NOTE: You can clone a specific branch using a similar method

```
git clone -b my-branch git@github.com:user/myproject.git
```

However, it still fetches all branches. To fetch only one branch...

```
mkdir $REMOTE_REPO-$BRANCH
cd $REMOTE_REPO-$BRANCH
git init
git remote add -t $BRANCH -f origin $REMOTE_REPO
git checkout $BRANCH
```

## Aliases

I created the following aliases to make things easier

```
git config --global alias.co checkout
git config --global alias.po "push origin"
git config --global alias.m merge
git config --global alias.nosslclone "-c http.sslVerify=false clone"
```

Get a list of current aliases with the following command

```
git config --get-regexp alias
```

-30-
