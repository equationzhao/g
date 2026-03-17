import pandas as pd
import streamlit as st

st.set_page_config(
    page_title="g - Advanced ls Command Alternative",
    page_icon="📂",
    # layout="wide",
    initial_sidebar_state="expanded"
)

st.markdown("""
<style>
    /* Base styles */
    body {
        font-family: 'Inter', sans-serif;
        line-height: 1.6;
    }
    
    /* Title and headers */
    .main-title {
        font-size: 3.5rem;
        font-weight: 800;
        color: #1E88E5;
        margin-bottom: 1.5rem;
        text-align: center;
        text-shadow: 0px 2px 3px rgba(0,0,0,0.1);
    }
    .section-header {
        font-size: 2.2rem;
        font-weight: 700;
        color: #0D47A1;
        border-bottom: 3px solid #90CAF9;
        padding-bottom: 0.7rem;
        margin-top: 2.5rem;
        margin-bottom: 1.5rem;
    }
    
    /* Cards and containers */
    .feature-title {
        font-weight: 700;
        font-size: 1.35rem;
        color: #1565C0;
        margin-bottom: 12px;
    }
    
    /* Code elements */
    .code-header {
        font-weight: 600;
        margin-top: 1.2rem;
        margin-bottom: 0.5rem;
        font-size: 1.1rem;
        color: #333;
    }
    .stCodeBlock {
        border-radius: 8px !important;
        margin-bottom: 1.5rem !important;
    }
    
    /* Images */
    .screenshot {
        border-radius: 12px;
        box-shadow: 0 8px 16px rgba(0,0,0,0.15);
        margin-bottom: 2rem;
    }
    
    /* Badges */
    .badge-container {
        display: flex;
        justify-content: center;
        flex-wrap: wrap;
        gap: 12px;
        margin-bottom: 25px;
        padding: 10px;
        background: rgba(240,240,240,0.5);
        border-radius: 8px;
    }
    .badge-container img {
        height: 28px;
        transition: transform 0.15s;
    }
    .badge-container img:hover {
        transform: scale(1.05);
    }
    
    /* Tables */
    .dataframe {
        font-size: 0.95rem !important;
    }
    .dataframe th {
        background-color: #E3F2FD !important;
        color: #0D47A1 !important;
    }
    
    /* Custom tabs styling */
    .stTabs [data-baseweb="tab-list"] {
        gap: 10px;
    }
    .stTabs [data-baseweb="tab"] {
        border-radius: 4px 4px 0 0;
        padding: 10px 16px;
        font-weight: 500;
    }
    
    /* Sidebar improvements */
    [data-testid="stSidebar"] [data-testid="stMarkdownContainer"] a {
        color: #1E88E5;
        text-decoration: none;
        transition: color 0.2s;
    }
    [data-testid="stSidebar"] [data-testid="stMarkdownContainer"] a:hover {
        color: #0D47A1;
        text-decoration: underline;
    }
    
    /* Footer */
    .footer {
        margin-top: 3rem;
        text-align: center;
        padding: 1.5rem;
        background-color: #f5f7f9;
        border-radius: 8px;
        color: #555;
    }
    
    /* Expander refinements */
    .streamlit-expanderHeader {
        font-weight: 600;
        font-size: 1.1rem;
        color: #1565C0;
    }
    
    /* Info block */
    .info-block {
        background-color: #E3F2FD;
        padding: 15px;
        border-radius: 8px;
        border-left: 5px solid #1E88E5;
        margin-bottom: 20px;
    }
</style>
""", unsafe_allow_html=True)

with st.sidebar:
    st.title("📋 Navigation")
    st.markdown("---")
    
    st.markdown("### 📑 Content")
    st.markdown("- [Key Features](#key-features)")
    st.markdown("- [Screenshots](#screenshots)")
    st.markdown("- [Usage](#usage)")
    st.markdown("- [Installation Guide](#installation-guide)")
    st.markdown("- [Shell Integration](#shell-integration)")
    st.markdown("- [Project Comparison](#project-comparison)")
    
    st.markdown("---")
    st.markdown("### 🔗 Quick Links")
    st.markdown("- [GitHub Repository](https://github.com/Equationzhao/g)")
    st.markdown("- [Report Issues](https://github.com/Equationzhao/g/issues)")
    st.markdown("- [Theme Documentation](https://github.com/Equationzhao/g/blob/master/docs/Theme.md)")
    st.markdown("- [Manual](https://github.com/Equationzhao/g/blob/master/docs/man.md)")
    
    st.markdown("---")
    st.caption("G - Version 0.30.0")

st.markdown('<div class="main-title">🌈 g </div>', unsafe_allow_html=True)

st.markdown('<div class="badge-container">'
            '<img src="https://img.shields.io/github/stars/Equationzhao/g" alt="Stars">'
            '<img src="https://img.shields.io/github/forks/Equationzhao/g" alt="Forks">'
            '<img src="https://img.shields.io/github/issues/Equationzhao/g" alt="Issues">'
            '<img src="https://img.shields.io/github/license/Equationzhao/g" alt="License">'
            '</div>', unsafe_allow_html=True)

st.markdown("""
**g** is a feature-rich, customizable, and cross-platform `ls` alternative. It provides enhanced visuals with type-specific icons, various layout options, and Git status integration.
""")

st.markdown('<div class="section-header" id="key-features">Key Features</div>', unsafe_allow_html=True)

col1, col2 = st.columns(2)

with col1:
    with st.container():
        st.markdown('<div class="feature-title">🎨 Customizable Display</div>', unsafe_allow_html=True)
        st.markdown("Icons and colors specific to file types, easy to customize")
        st.markdown('</div>', unsafe_allow_html=True)
        st.markdown('</div>', unsafe_allow_html=True)

    st.markdown("")  # Spacing

    with st.container():
        st.markdown('<div class="feature-title">🔀 Multiple Layouts</div>', unsafe_allow_html=True)
        st.markdown("Choose from grid, across, byline, zero, comma, table, json, markdown, and tree layouts")
        st.markdown('</div>', unsafe_allow_html=True)
        st.markdown('</div>', unsafe_allow_html=True)

    st.markdown("")  # Spacing

    with st.container():
        st.markdown('<div class="feature-title">🌐 Git Integration</div>', unsafe_allow_html=True)
        st.markdown("View file git-status/repo-status/repo-branch directly in your listings")
        st.markdown('</div>', unsafe_allow_html=True)
        st.markdown('</div>', unsafe_allow_html=True)

with col2:
    with st.container():
        st.markdown('<div class="feature-title">🔄 Advanced Sorting</div>', unsafe_allow_html=True)
        st.markdown("Highly customizable sorting options like version-sort")
        st.markdown('</div>', unsafe_allow_html=True)
        st.markdown('</div>', unsafe_allow_html=True)

    st.markdown("")  # Spacing

    with st.container():
        st.markdown('<div class="feature-title">💻 Cross-Platform</div>', unsafe_allow_html=True)
        st.markdown("Works seamlessly on Linux, Windows, and MacOS")
        st.markdown('</div>', unsafe_allow_html=True)
        st.markdown('</div>', unsafe_allow_html=True)

    st.markdown("")  # Spacing

    with st.container():
        st.markdown('<div class="feature-title">🔍 Fuzzy Path Matching</div>', unsafe_allow_html=True)
        st.markdown("zoxide and fzf like fuzzy path matching")
        st.markdown('</div>', unsafe_allow_html=True)
        st.markdown('</div>', unsafe_allow_html=True)

st.markdown('<div class="section-header" id="screenshots">Screenshots</div>', unsafe_allow_html=True)
st.markdown('<div class="screenshot-container">', unsafe_allow_html=True)
st.image("https://raw.githubusercontent.com/Equationzhao/g/master/asset/screenshot_3.png", 
         caption="g command interface with icons and features visible", 
         use_column_width=True, 
         output_format="PNG")
st.markdown('</div>', unsafe_allow_html=True)

st.markdown('<div class="section-header" id="usage">Usage</div>', unsafe_allow_html=True)

st.markdown('<div class="code-header">Basic usage:</div>', unsafe_allow_html=True)
st.code("g path(s)", language="bash")

st.markdown('<div class="code-header">Display icons and long format:</div>', unsafe_allow_html=True)
st.code("g --icon --long path(s)", language="bash")

st.markdown('<div class="code-header">Tree layout:</div>', unsafe_allow_html=True)
st.code("g --tree --long path(s)", language="bash")

more_options = st.expander("View more options")
with more_options:
    st.markdown("""
    ### Layout Options
    - `--grid`: Grid layout (default)
    - `--tree`: Tree layout
    - `--across`: Across layout
    - `--oneline`: One line per entry
    - `--table`: Table layout
    - `--json`: JSON output
    - `--markdown`: Markdown table

    ### Display Options
    - `--icon`: Show file type icons
    - `--long`, `-l`: Long format with details
    - `--all`, `-a`: Show hidden files
    - `--git`: Show git status
    - `--mime`: Show MIME types
    - `--hyperlink`: Enable file hyperlinks

    For all options, see the [manual](https://github.com/Equationzhao/g/blob/master/docs/man.md)
    """)
    
theme = st.expander("Custom Themes", expanded=False)
with theme:
    st.markdown("""
    g supports rich theme customization.
    
    For detailed information, see the [Theme Documentation](https://github.com/Equationzhao/g/blob/master/docs/Theme.md).
    
    Basic usage:
    ```bash
    g --theme mytheme.json  # Use custom theme
    ```
    """)

st.markdown('<div class="section-header" id="installation-guide">Installation Guide</div>', unsafe_allow_html=True)

install_tabs = st.tabs(["📦 Package Managers", "⬇️ Pre-built Binaries", "🛠️ Build from Source"])

with install_tabs[0]:
    col1, col2 = st.columns(2)
    
    with col1:
        st.subheader("Arch Linux")
        st.code("yay -S g-ls", language="bash")
        
        st.subheader("Homebrew")
        st.code("brew install g-ls", language="bash")
        
        st.subheader("MacPort")
        st.code("sudo port install g-ls", language="bash")
    
    with col2:
        st.subheader("Windows (Scoop)")
        st.code("scoop install https://raw.githubusercontent.com/Equationzhao/g/master/scoop/g.json", language="powershell")
        
        st.subheader("WinGet")
        st.info("In development, see [#119](https://github.com/Equationzhao/g/issues/119)")

with install_tabs[1]:
    col1, col2 = st.columns(2)
    
    with col1:
        st.subheader("Install Script")
        st.code('bash -c "$(curl -fsSLk https://raw.githubusercontent.com/Equationzhao/g/master/script/install.sh)"', language="bash")
        
        st.subheader("Uninstall")
        st.code('curl -fsSLk https://raw.githubusercontent.com/Equationzhao/g/master/script/install.sh | bash /dev/stdin -r', language="bash")
    
    with col2:
        st.subheader("Download Packages")
        st.markdown("""
        1. Download from [GitHub releases page](https://github.com/Equationzhao/g/releases)
        2. Extract and add the executable to your PATH
        
        #### Debian/Ubuntu
        ```bash
        sudo dpkg -i g_$version_$arch.deb
        ```
        """)

with install_tabs[2]:
    st.markdown("""
    ### Requirements
    - Go version >= 1.26.1
    
    ### Install Latest Version
    ```bash
    go install -ldflags="-s -w" github.com/Equationzhao/g@latest
    ```
    
    ### Build from Source
    ```bash
    git clone github.com/Equationzhao/g
    cd g
    go build -ldflags="-s -w"
    # then add the executable file to your PATH
    ```
    """)

st.markdown("### Recommended Terminals")

term_cols = st.columns(3)

with term_cols[0]:
    st.markdown("#### macOS")
    st.markdown("- [Iterm2](https://iterm2.com/)")
    st.markdown("- [Warp](https://www.warp.dev)")

with term_cols[1]:
    st.markdown("#### Windows")
    st.markdown("- [Windows Terminal](https://github.com/microsoft/terminal)")

with term_cols[2]:
    st.markdown("#### Cross-platform")
    st.markdown("- [Hyper](https://hyper.is/)")
    st.markdown("- [WezTerm](https://wezfurlong.org/wezterm/)")

st.markdown('<div class="section-header" id="shell-integration">Shell Integration</div>', unsafe_allow_html=True)

shell_tabs = st.tabs(["🔄 Command Completion", "🔠 Shell Aliases"])

with shell_tabs[0]:
    st.markdown('<div class="info-block">', unsafe_allow_html=True)
    st.info("If you install `g` through brew or the install script, the completion is usually installed already.")
    st.markdown('</div>', unsafe_allow_html=True)
    
    completion_tabs = st.tabs(["zsh", "bash", "fish"])
    
    with completion_tabs[0]:
        st.code("""
# 1. Download completion script
wget https://raw.githubusercontent.com/Equationzhao/g/master/completions/zsh/_g

# 2. Install to completion directory
mv _g ~/.zsh/completions  # or any directory in your $FPATH

# 3. Make sure these commands are in your .zshrc
autoload -Uz compinit
compinit
        """, language="bash")
    
    with completion_tabs[1]:
        st.code("""
# 1. Download completion script
wget https://raw.githubusercontent.com/Equationzhao/g/master/completions/bash/g-completion.bash

# 2. Add to .bashrc
echo "source /path/to/g-completion.bash" >> ~/.bashrc

# 3. Reload configuration
source ~/.bashrc
        """, language="bash")
    
    with completion_tabs[2]:
        st.code("""
# 1. Download completion script
wget https://raw.githubusercontent.com/Equationzhao/g/master/completions/fish/g.fish

# 2. Install to completion directory
mv g.fish ~/.config/fish/completions

# 3. Reload configuration
source ~/.config/fish/config.fish
        """, language="fish")

with shell_tabs[1]:
    alias_tabs = st.tabs(["bash", "zsh", "fish", "powershell", "nushell"])
    
    with alias_tabs[0]:
        st.code("""
# Add to .bashrc
eval "$(g --init bash)"

# Reload configuration
source ~/.bashrc
        """, language="bash")
    
    with alias_tabs[1]:
        st.code("""
# Add to .zshrc
eval "$(g --init zsh)"

# Reload configuration
source ~/.zshrc
        """, language="zsh")
    
    with alias_tabs[2]:
        st.code("""
# Add to fish configuration
g --init fish | source

# Reload configuration
source ~/.config/fish/config.fish
        """, language="fish")
    
    with alias_tabs[3]:
        st.code("""
# Add to your profile
Invoke-Expression (& { (g --init powershell | Out-String) })

# Find your profile path
echo $profile
        """, language="powershell")
    
    with alias_tabs[4]:
        st.code("""
# Add to $nu.env-path
^g --init nushell | save -f ~/.g.nu

# Add to $nu.config-path
source ~/.g.nu
        """, language="nu")

st.markdown('<div class="section-header" id="project-comparison">Project Comparison</div>', unsafe_allow_html=True)

st.markdown("""
g is highly inspired by the following excellent projects:
- [exa](https://github.com/ogham/exa) / [eza](https://github.com/eza-community/eza)
- [lsd](https://github.com/lsd-rs/lsd)
- [ls-go](https://github.com/acarl005/ls-go)
""")

comparison_data = {
    "Feature": ["Display Modes", "Unique Features", "Performance"],
    "eza": ["oneline, grid, across, tree, recurse", 
            "-Z: list security context, -@: list extended attributes and sizes", 
            "better"],
    "g": ["oneline, grid, across, zero, comma, table, json, markdown, tree, recurse", 
          "--mime: list mime types, --charset: list character sets", 
          "slower"]
}

df = pd.DataFrame(comparison_data)
st.table(df)

st.markdown("---")
st.markdown("""
<div class="footer">
    <p>g is an open-source project with MIT license</p>
    <p>© 2023-2025 Equationzhao</p>
    <p><a href="https://github.com/Equationzhao/g/stargazers">⭐ Star on GitHub</a> if you find this project useful!</p>
</div>
""", unsafe_allow_html=True)
