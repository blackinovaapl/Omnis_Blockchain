<h1 align="center">✨ Omnis Blockchain ✨</h1>

<p align="center">
  The <strong>Omnis Blockchain</strong> is a custom Layer-1 blockchain built from the ground up using the <strong>Cosmos SDK</strong>.<br>
  Our mission is to provide a robust and flexible foundation for decentralized applications by introducing the <strong>custom OMS-20 Token Standard</strong>,<br>
  a powerful framework for creating and managing fungible assets.
</p>

<hr>

<h2>🚀 Key Features</h2>
<ul>
  <li><strong>Custom OMS-20 Token Standard</strong> – A purpose-built, extensible standard for fungible tokens, designed to be flexible and secure.</li>
  <li><strong>Built on Cosmos SDK</strong> – Leverages the battle-tested framework of the Cosmos SDK for scalability, security, and interoperability.</li>
  <li><strong>Protobuf-Based APIs</strong> – Core logic and data structures are defined using Protocol Buffers for language-agnostic integration.</li>
  <li><strong>Consensus Powered by Tendermint BFT</strong> – Rapid finality and high throughput using the robust Tendermint BFT engine.</li>
  <li><strong><code>omnisd</code> CLI Tool</strong> – Easily interact with the blockchain, manage keys, and execute transactions.</li>
</ul>

<hr>

<h2>🛠️ Technologies & Skills</h2>
<ul>
  <li><strong>Go</strong> – Core blockchain programming language.</li>
  <li><strong>Cosmos SDK</strong> – Blockchain framework foundation.</li>
  <li><strong>Protobuf</strong> – Defines data structures and APIs.</li>
  <li><strong>Buf</strong> – Linting and code generation tool for Protobuf.</li>
  <li><strong>Tendermint BFT</strong> – Consensus mechanism.</li>
  <li><strong>YAML & Shell</strong> – Configuration and automation tools.</li>
</ul>

<hr>

<h2>🏁 Getting Started</h2>
<p><em>Prerequisites</em>: Go (v1.24.4+), Buf, Ignite</p>

<ol>
  <li>
    <strong>Build the Blockchain</strong>
    <pre><code>cd ~/omnis
ignite chain build</code></pre>
  </li>
  <li>
    <strong>Start a Local Network</strong>
    <pre><code>ignite chain serve</code></pre>
  </li>
  <li>
    <strong>Interact with the Blockchain</strong>
    <pre><code># Create a wallet key
omnisd keys add mywallet

# Check account balance
omnisd query bank balances &lt;your-wallet-address&gt;

# Create a new OMS-20 Token
omnisd tx token create-token "My Awesome Token" "MAT" "1000000000" --from mywallet</code></pre>
  </li>
</ol>

<hr>

<h2>📁 Project Structure</h2>
<pre>
omnis/
├── app/                  # Application-level configurations
├── proto/                # Protobuf definitions
│   ├── omnis/
│   │   ├── omnis/
│   │   └── token/
│   └── third_party/      # Imported Protobuf definitions
├── x/                    # Custom modules
│   ├── omnis/
│   └── token/            # OMS-20 module logic
├── gen/go/               # Generated Go code from Protobuf
├── buf.gen.yaml          # Buf code generation config
├── buf.yaml              # Buf linting and workspace config
├── go.mod                # Go module file
└── README.md             # This document
</pre>

<hr>

<h2>🗺️ Future Roadmap</h2>
<ul>
  <li>Full OMS-20 Implementation – Token transfers, permissions, and more.</li>
  <li>Omnis Block Explorer – Web-based network activity visualization.</li>
  <li>Testnet Deployment – Public multi-node testnet.</li>
  <li>Custom Consensus – Long-term research and implementation.</li>
</ul>

<hr>

<h2>🙏 Acknowledgements</h2>
<ul>
  <li><a href="https://github.com/cosmos/cosmos-sdk">Cosmos SDK</a></li>
  <li><a href="https://github.com/ignite">Ignite</a></li>
  <li><a href="https://buf.build">Buf</a></li>
</ul>
