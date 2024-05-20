#!/bin/bash

# Update and install Vim
sudo apt update
sudo apt install -y vim

# Install vim-plug
curl -fLo ~/.vim/autoload/plug.vim --create-dirs \
    https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim

# Create .vimrc with Copilot plugin
cat <<EOL > ~/.vimrc
call plug#begin('~/.vim/plugged')

" GitHub Copilot Vim plugin
Plug 'github/copilot.vim'

call plug#end()

" Enable Copilot
let g:copilot_no_tab_map = v:true
imap <silent><script><expr> <C-J> copilot#Accept("\<CR>")
EOL

# Install Vim plugins
vim +PlugInstall +qall
