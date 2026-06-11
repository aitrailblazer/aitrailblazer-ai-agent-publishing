package main

import (
	"net/http"
	"os"
	"strings"
)

const publishingDemoLandingHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>AITrailblazer AI Agent Publishing</title>
  <style>
    :root{color-scheme:dark;--text:#f7fbff;--muted:#b8c5d8;--green:#61d394;--blue:#70d7ff;--gold:#d8b35f}
    *{box-sizing:border-box}
    html,body{margin:0;width:100%;height:100%;overflow:hidden;background:#05070b;color:var(--text);font:15px/1.35 Inter,ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;-webkit-font-smoothing:antialiased;text-rendering:geometricPrecision}
    body{display:grid;place-items:center;padding:10px;background:radial-gradient(circle at 78% 8%,rgba(97,211,148,.18),transparent 28%),radial-gradient(circle at 12% 12%,rgba(112,215,255,.14),transparent 28%),linear-gradient(135deg,#04070d,#08111b 54%,#06110c)}
    main{width:min(1240px,100%);height:calc(100dvh - 20px);display:grid;grid-template-columns:1fr .92fr;gap:14px;align-items:stretch;overflow:hidden}
    .hero,.panel{min-height:0;border:1px solid rgba(255,255,255,.14);border-radius:22px;background:linear-gradient(180deg,rgba(18,28,43,.96),rgba(7,12,20,.98));box-shadow:0 30px 90px rgba(0,0,0,.42),inset 0 1px 0 rgba(255,255,255,.08)}
    .hero{position:relative;overflow:hidden;padding:22px;display:flex;flex-direction:column;justify-content:space-between;border-color:rgba(216,179,95,.46)}
    .hero::after{content:"";position:absolute;right:-130px;top:-190px;width:440px;height:440px;border-radius:50%;background:radial-gradient(circle,rgba(216,179,95,.22),rgba(97,211,148,.12),transparent 65%);animation:float 7s ease-in-out infinite}
    .brand{position:relative;z-index:1;display:flex;align-items:center;gap:12px}
    .logo{width:54px;height:54px;border-radius:16px;object-fit:cover;box-shadow:0 0 0 1px rgba(255,255,255,.18),0 0 34px rgba(216,179,95,.34)}
    .eyebrow{color:var(--gold);text-transform:uppercase;letter-spacing:.11em;font-weight:950;font-size:11px}
    h1{position:relative;z-index:1;margin:14px 0 8px;font-size:clamp(38px,4.2vw,52px);line-height:.94;letter-spacing:0}
    .lead{position:relative;z-index:1;color:#eef6ff;font-size:17px;line-height:1.2;font-weight:820;max-width:760px}
    .voiceover{position:relative;z-index:1;margin-top:10px;border:1px solid rgba(112,215,255,.24);border-radius:16px;background:rgba(112,215,255,.07);padding:10px;color:#dff6ff;font-size:13px;line-height:1.24;font-weight:780;max-width:780px}
    .voiceover b{color:#f8d998}
    .qualify{position:relative;z-index:1;margin-top:8px;display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:5px;max-width:780px}
    .check{display:flex;align-items:flex-start;gap:6px;border:1px solid rgba(97,211,148,.24);border-radius:11px;background:rgba(97,211,148,.07);padding:6px;color:#eafcff;font-size:10.2px;font-weight:820;line-height:1.08}
    .box{display:grid;place-items:center;flex:0 0 17px;width:17px;height:17px;border-radius:5px;background:linear-gradient(135deg,#61d394,#70d7ff);color:#041018;font-size:12px;font-weight:950;box-shadow:0 0 20px rgba(97,211,148,.22)}
    .actions{position:relative;z-index:1;display:flex;gap:9px;flex-wrap:wrap;margin-top:10px}
    a.button{display:inline-flex;align-items:center;gap:10px;text-decoration:none;border-radius:999px;padding:10px 14px;color:#06110c;background:linear-gradient(135deg,#81e6a5,#70d7ff);font-weight:950;box-shadow:0 20px 48px rgba(97,211,148,.22);transition:transform .16s ease}
    a.button:hover{transform:translateY(-2px)}a.button.secondary{color:#f7fbff;background:rgba(112,215,255,.1);border:1px solid rgba(112,215,255,.32)}
    .trip{display:inline-flex;border:1px solid rgba(216,179,95,.5);border-radius:999px;padding:8px 10px;color:#f8d998;background:rgba(216,179,95,.08);font:900 12px "SF Mono",ui-monospace,monospace;letter-spacing:.04em}
    .panel{padding:16px;display:grid;grid-template-rows:auto 1fr;gap:10px;overflow:hidden}
    .panel h2{margin:0;color:var(--green);font-size:21px}.panel p{font-size:13px;font-weight:760;line-height:1.28}
    .steps{display:grid;gap:7px;min-height:0}
    .summary{display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:6px}
    .summary div{border:1px solid rgba(255,255,255,.12);border-radius:12px;background:rgba(255,255,255,.045);padding:7px}
    .summary b{display:block;color:#fff;font-size:13px}.summary span{display:block;color:var(--muted);font-size:11.5px;margin-top:3px;line-height:1.2}
    .step{display:grid;grid-template-columns:34px 1fr;gap:9px;align-items:center;border:1px solid rgba(255,255,255,.12);border-radius:14px;padding:8px;background:rgba(255,255,255,.045);animation:reveal .7s both}
    .step:nth-child(2){animation-delay:.12s}.step:nth-child(3){animation-delay:.24s}.step:nth-child(4){animation-delay:.36s}.step:nth-child(5){animation-delay:.48s}
    .num{display:grid;place-items:center;width:30px;height:30px;border-radius:10px;background:rgba(97,211,148,.16);color:#8ff2b4;font-weight:950;border:1px solid rgba(97,211,148,.34)}
    .step b{display:block;color:#fff;font-size:13.5px}.step span{display:block;color:var(--muted);font-size:12px;margin-top:2px;line-height:1.18}
    .diagram{height:124px;border:1px solid rgba(112,215,255,.24);border-radius:16px;background:radial-gradient(circle at 50% 50%,rgba(97,211,148,.13),transparent 42%),rgba(0,0,0,.18);display:grid;place-items:center;overflow:hidden}
    svg{width:100%;height:100%}.node{fill:#111b2a;stroke:#61d394;stroke-width:3;filter:drop-shadow(0 0 12px rgba(97,211,148,.24))}.node.cold{stroke:#70d7ff}.wire{fill:none;stroke:#70d7ff;stroke-width:4;stroke-dasharray:12 10;animation:dash 1.3s linear infinite}.txt{fill:#f7fbff;font-size:19px;font-weight:900}.sub{fill:#b8c5d8;font-size:12px;font-weight:800}
    footer{position:fixed;right:18px;bottom:12px;color:rgba(247,251,255,.42);font-size:11px;text-transform:uppercase;letter-spacing:.08em}
    @keyframes dash{to{stroke-dashoffset:-44}}@keyframes float{50%{transform:translateY(20px)}}@keyframes reveal{from{opacity:0;transform:translateY(12px)}to{opacity:1;transform:translateY(0)}}
    @media(max-height:760px) and (min-width:981px){
      body{padding:8px}
      main{height:calc(100dvh - 16px);gap:12px}
      .hero{padding:18px}.panel{padding:14px}
      .logo{width:48px;height:48px}.eyebrow{font-size:10px}
      h1{font-size:clamp(36px,3.9vw,48px);margin:12px 0 7px}
      .lead{font-size:15.5px;line-height:1.15}
      .voiceover{font-size:12px;line-height:1.16;padding:8px;margin-top:8px}
      .check{font-size:10px;line-height:1.08;padding:5px 6px}.box{width:15px;height:15px;flex-basis:15px;font-size:10px}
      .actions{margin-top:8px}.button{font-size:13px}.trip{font-size:11px}
      .panel h2{font-size:20px}.panel p{font-size:12px}
      .diagram{height:112px}.summary span{font-size:10.5px}.summary b{font-size:12px}
      .step{padding:7px}.step b{font-size:12.5px}.step span{font-size:11px}.num{width:28px;height:28px}
    }
    @media(max-width:980px){html,body{overflow:auto}main{height:auto;grid-template-columns:1fr}.hero{min-height:auto}body{display:block}.lead{font-size:17px}h1{font-size:44px}.qualify{grid-template-columns:1fr}}
  </style>
</head>
<body>
<main>
  <section class="hero">
    <div>
      <div class="brand">
        <img class="logo" src="/img/DeltaSignalIcon-256.png" alt="DeltaSignal logo">
        <div><div class="eyebrow">Google Cloud Rapid Agent Hackathon · MongoDB Track</div><strong>AITrailblazer AI Agent Publishing</strong></div>
      </div>
      <h1>Turn a publication archive into agent memory.</h1>
      <p class="lead">AITrailblazer AI Agent Publishing converts published research into a reusable agent workflow: archive context, TripCode identity, River memory, MongoDB-backed state, Gemini-ready synthesis, and a rendered user packet.</p>
      <div class="voiceover"><b>Opening voiceover:</b> This demo shows a publishing agent that starts with research content, assigns a stable TripCode, reconstructs the surrounding River, remembers the session, and returns both raw JSON proof and a polished HTML packet that another agent or reader can reuse.</div>
      <div class="qualify" aria-label="Submission qualification checklist">
        <div class="check"><span class="box">✓</span><span>Functional running application, not a slide-only concept.</span></div>
        <div class="check"><span class="box">✓</span><span>Google Cloud runtime surface with browser and API demo paths.</span></div>
        <div class="check"><span class="box">✓</span><span>Gemini-ready synthesis loop with bounded evidence packets.</span></div>
        <div class="check"><span class="box">✓</span><span>MongoDB-oriented memory model for archive, TripCode, and River state.</span></div>
        <div class="check"><span class="box">✓</span><span>Repeatable 3-minute demo with replayable commands and UI actions.</span></div>
        <div class="check"><span class="box">✓</span><span>Cost-aware usage ledger and protected demo-key access.</span></div>
      </div>
      <div class="actions">
        <a class="button" href="/demo/run?autoplay=1">Run 3 minute demo →</a>
        <a class="button secondary" href="/demo/run">Open manual controls</a>
        <span class="trip">HUT-RIVER-001</span>
      </div>
    </div>
  </section>
  <section class="panel">
    <div>
      <h2>What the browser demo proves</h2>
      <p>Start with a publication archive, run the deployed Cloud Run agent, inspect machine-readable JSON, then render the same response as a user-facing research packet.</p>
    </div>
    <div class="steps">
      <div class="diagram"><svg viewBox="0 0 720 210"><path class="wire" d="M130 72 H290M430 72 H590M360 110 V158"/><rect class="node" x="42" y="34" width="120" height="76" rx="20"/><text class="txt" x="73" y="70">Archive</text><text class="sub" x="66" y="92">articles</text><rect class="node" x="300" y="34" width="120" height="76" rx="20"/><text class="txt" x="322" y="70">Agent</text><text class="sub" x="315" y="92">Cloud Run</text><rect class="node cold" x="558" y="34" width="120" height="76" rx="20"/><text class="txt" x="590" y="70">MCP</text><text class="sub" x="580" y="92">MongoDB</text><rect class="node" x="270" y="134" width="180" height="58" rx="18"/><text class="txt" x="305" y="169">HTML Packet</text></svg></div>
      <div class="summary">
        <div><b>Problem</b><span>Research archives are valuable, but passive prose is hard for agents to execute.</span></div>
        <div><b>Agent loop</b><span>Discover, invoke, verify, remember, render, and reuse the packet.</span></div>
        <div><b>Demo proof</b><span>The same route supports manual testing and the 3-minute video sequence.</span></div>
      </div>
      <div class="step"><div class="num">1</div><div><b>Archive brief</b><span>Publication inventory becomes article objects, TripCodes, Rivers, and reusable reader prompts.</span></div></div>
      <div class="step"><div class="num">2</div><div><b>TripCode resolve</b><span>HUT-RIVER-001 becomes structured River memory with MongoDB fit and execution trace.</span></div></div>
      <div class="step"><div class="num">3</div><div><b>Gemini-ready packet</b><span>The agent returns evidence boundaries, memory follow-up, usage visibility, and synthesis inputs.</span></div></div>
      <div class="step"><div class="num">4</div><div><b>Rendered output</b><span>The latest JSON response becomes a polished HTML report for users, agents, and reviewers.</span></div></div>
    </div>
  </section>
</main>
<footer>© 2026 AITrailblazer</footer>
</body>
</html>`

const publishingDemoUIHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>AITrailblazer AI Agent Publishing Demo</title>
  <style>
    :root{color-scheme:dark;--bg:#05070b;--panel:#101722;--line:rgba(255,255,255,.14);--text:#f7fbff;--muted:#b8c5d8;--green:#61d394;--blue:#70d7ff;--gold:#d8b35f;--amber:#f8d998;--red:#ff6874;--mono:"SF Mono","JetBrains Mono","Cascadia Code",ui-monospace,SFMono-Regular,Menlo,Monaco,Consolas,monospace}
    *{box-sizing:border-box}
    html,body{margin:0;width:100%;height:100%;overflow:hidden;background:#05070b;color:var(--text);font:15px/1.38 Inter,ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;-webkit-font-smoothing:antialiased;text-rendering:geometricPrecision}
    body{background:radial-gradient(circle at 76% 8%,rgba(97,211,148,.16),transparent 26%),radial-gradient(circle at 10% 14%,rgba(112,215,255,.11),transparent 26%),linear-gradient(135deg,#04070d,#08111b 52%,#06110c)}
    publishing-demo{display:block;width:100vw;height:100vh;overflow:hidden}
    main{height:100dvh;max-width:1340px;margin:0 auto;padding:12px;display:grid;grid-template-rows:auto minmax(0,1fr) auto;gap:10px;overflow:hidden}
    header,.panel,.captionBar{border:1px solid var(--line);border-radius:16px;background:linear-gradient(180deg,rgba(18,28,43,.98),rgba(7,12,20,.98));box-shadow:0 20px 58px rgba(0,0,0,.38),inset 0 1px 0 rgba(255,255,255,.08)}
    header{position:relative;display:grid;grid-template-columns:minmax(330px,1fr) minmax(440px,620px) auto;align-items:center;gap:12px;padding:10px 12px;border-color:rgba(216,179,95,.45);overflow:hidden}
    header::after{content:"";position:absolute;right:-70px;top:-120px;width:300px;height:300px;border-radius:50%;background:radial-gradient(circle,rgba(216,179,95,.25),rgba(97,211,148,.12),transparent 64%)}
    .brand{position:relative;z-index:1;display:flex;align-items:center;gap:12px;min-width:0}
    .brand a{display:flex;align-items:center;gap:12px;color:inherit;text-decoration:none;min-width:0}
    .logo{width:48px;height:48px;border-radius:14px;object-fit:cover;box-shadow:0 0 0 1px rgba(255,255,255,.18),0 0 28px rgba(216,179,95,.34)}
    .eyebrow{color:var(--gold);text-transform:uppercase;letter-spacing:.12em;font-weight:950;font-size:12px}
    h1{font-size:clamp(28px,3.1vw,43px);line-height:.98;margin:3px 0 0;letter-spacing:0}
    .status{position:relative;z-index:1;display:inline-flex;align-items:center;gap:8px;border:1px solid rgba(97,211,148,.45);border-radius:999px;padding:9px 13px;color:#8ff2b4;background:rgba(97,211,148,.1);font-weight:950;white-space:nowrap}
    .dot{width:9px;height:9px;border-radius:999px;background:var(--green);box-shadow:0 0 18px rgba(97,211,148,.7)}
    .grid{min-height:0;display:grid;grid-template-columns:38.2fr 61.8fr;gap:10px;overflow:hidden}
    .panel{min-height:0;overflow:hidden;padding:12px}
    .control{display:grid;grid-template-rows:auto auto auto minmax(0,1fr);gap:9px;border-color:rgba(216,179,95,.28)}
    .workspace{display:grid;grid-template-rows:auto minmax(0,1fr);gap:10px;border-color:rgba(112,215,255,.62);background:radial-gradient(circle at 100% 0%,rgba(112,215,255,.18),transparent 32%),linear-gradient(180deg,rgba(6,31,57,.99),rgba(1,12,27,.995));box-shadow:0 24px 70px rgba(0,0,0,.48),0 0 0 2px rgba(112,215,255,.15)}
    .panelHead{display:flex;align-items:flex-start;justify-content:space-between;gap:10px}
    h2{margin:0;color:var(--green);font-size:20px;line-height:1.05}
    p{margin:0;color:var(--muted)}
    .pill{display:inline-flex;align-items:center;border:1px solid rgba(97,211,148,.36);border-radius:999px;padding:6px 9px;color:#8ff2b4;background:rgba(97,211,148,.1);font-size:11px;font-weight:950}
    .metrics{display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:7px}
    .metric,.info{border:1px solid rgba(255,255,255,.12);border-radius:13px;background:linear-gradient(180deg,rgba(255,255,255,.065),rgba(255,255,255,.025));padding:9px}
    .metric strong{display:block;color:#8ff2b4;font-size:21px}.metric span,.info span{display:block;color:var(--muted);font-size:11px;margin-top:3px}
    .flow{display:grid;grid-template-columns:repeat(5,minmax(0,1fr));gap:5px;border:1px solid rgba(255,255,255,.1);border-radius:13px;padding:5px;background:rgba(0,0,0,.16)}
    .flow div{border:1px solid rgba(112,215,255,.18);border-radius:9px;padding:7px 4px;text-align:center;color:#ddecff;font-size:10px;font-weight:850}.flow b{display:block;color:#fff;font-size:11px}
    .inputs{position:relative;z-index:1;display:grid;grid-template-columns:1.35fr 1fr .9fr;gap:7px;align-items:end}
    label{display:block;color:var(--amber);font-weight:950;font-size:10px;text-transform:uppercase;letter-spacing:.09em;margin:0 0 4px}
    input{width:100%;height:34px;border:1px solid rgba(112,215,255,.32);background:#07101d;color:var(--text);border-radius:10px;padding:8px 10px;font:13px var(--mono)}
    .buttons{display:grid;grid-template-columns:repeat(4,minmax(0,1fr));gap:7px}
    button{position:relative;overflow:hidden;display:inline-flex;align-items:center;justify-content:center;gap:6px;border:1px solid rgba(97,211,148,.36);background:linear-gradient(135deg,rgba(97,211,148,.25),rgba(112,215,255,.13));color:var(--text);border-radius:10px;padding:8px 6px;font-weight:950;font-size:12px;cursor:pointer;box-shadow:0 10px 28px rgba(0,0,0,.22);transition:transform .14s ease,border-color .14s ease,box-shadow .14s ease}
    button:hover{border-color:var(--green);transform:translateY(-1px)}button.secondary{border-color:rgba(112,215,255,.32);background:rgba(112,215,255,.08)}
    button.running{border-color:var(--green);box-shadow:0 0 0 2px rgba(97,211,148,.16),0 0 34px rgba(97,211,148,.24);animation:pulse 1s ease-in-out infinite}
    button.running::after{content:"";position:absolute;inset:-45%;background:linear-gradient(90deg,transparent,rgba(255,255,255,.26),transparent);transform:translateX(-70%);animation:sweep 1.15s linear infinite}
    button.clicked{animation:clickpop .64s ease-out}
    button.clicked::before{content:"";position:absolute;inset:0;border-radius:inherit;background:radial-gradient(circle at 50% 50%,rgba(248,217,152,.48),transparent 62%);animation:ripple .62s ease-out;pointer-events:none}
    .btnIcon{position:relative;z-index:1;display:inline-grid;place-items:center;width:17px;height:17px;flex:0 0 17px;border-radius:6px;background:rgba(0,0,0,.2);box-shadow:inset 0 0 0 1px rgba(255,255,255,.12)}
    .btnIcon svg{width:12px;height:12px;stroke:#f8d998;stroke-width:2.1;fill:none;stroke-linecap:round;stroke-linejoin:round}
    .btnLabel{position:relative;z-index:1;white-space:nowrap}
    .explain{border:1px solid rgba(216,179,95,.28);border-radius:14px;background:linear-gradient(180deg,rgba(216,179,95,.08),rgba(255,255,255,.03));padding:11px;min-height:104px}
    .explain h3{margin:0 0 6px;color:#ffe4a3;font-size:15px}.explain p{font-size:17px;line-height:1.23;color:#eef6ff;font-weight:820}
    .outHead{display:flex;align-items:center;justify-content:space-between;gap:10px}.outHead p{font-size:13px}.tools{display:flex;gap:6px;align-items:center}.tools button{padding:7px 9px;font-size:11px}
    pre{white-space:pre-wrap;word-break:break-word;margin:0;height:100%;overflow:auto;border:1px solid rgba(112,215,255,.42);border-radius:16px;background:linear-gradient(180deg,#041424,#02070e);color:#eafcff;padding:17px;font:15px/1.58 var(--mono);box-shadow:inset 0 1px 0 rgba(255,255,255,.08),0 0 44px rgba(112,215,255,.08)}
    .resp-title{color:var(--green);font-weight:950}.resp-status{color:var(--blue);font-weight:950}.resp-error{color:var(--red);font-weight:950}
    .json-key{color:#7ddcff;font-weight:850}.json-string{color:#f6d489}.json-number{color:#caa8ff}.json-bool{color:#61d394;font-weight:850}.json-null{color:#ff94a8;font-weight:850}.json-punc{color:#7d8ca7}.json-brace{color:#d8f3ff;font-weight:900}.json-colon{color:#d6ad5c}
    .captionBar{min-height:84px;padding:12px 17px;border-color:rgba(248,217,152,.62);background:linear-gradient(135deg,#07090e,#10141c 48%,#22190b);box-shadow:0 -20px 64px rgba(0,0,0,.58),0 0 38px rgba(214,173,92,.16)}
    .captionText{display:block;color:#d3a63a;font:950 25px/1.08 "Arial Narrow","Roboto Condensed","Aptos Narrow","Inter Tight",Inter,ui-sans-serif,system-ui,sans-serif;max-height:3.4em;overflow:hidden;text-shadow:0 2px 0 rgba(0,0,0,.78)}
    .cursor{display:inline-block;width:9px;height:1em;margin-left:5px;vertical-align:-.12em;background:#d3a63a;animation:blink .78s steps(2,end) infinite}
    .trace{display:block;margin:0 0 13px;padding:12px;border:1px solid rgba(97,211,148,.24);border-radius:13px;background:linear-gradient(180deg,rgba(97,211,148,.09),rgba(112,215,255,.05))}
    .trace b{color:#ffe4a3}.trace .ok{color:#61d394}.route{color:#7ddcff}
    .rich{display:block}.richHero,.richSection{display:block;border:1px solid rgba(112,215,255,.26);border-radius:14px;padding:13px;margin-bottom:10px;background:rgba(255,255,255,.045)}.richHero{border-color:rgba(97,211,148,.32);background:linear-gradient(135deg,rgba(97,211,148,.13),rgba(112,215,255,.08))}
    .richKicker{display:block;color:#d8b35f;text-transform:uppercase;letter-spacing:.1em;font-weight:950;font-size:11px}.richTitle{display:block;color:#fff;font-size:25px;font-weight:950;line-height:1.04;margin-top:6px}.richSummary{display:block;color:#dce7f6;font-size:15px;line-height:1.34;margin-top:8px}
    .richGrid{display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:8px;margin-bottom:10px}.richCard{display:block;border:1px solid rgba(255,255,255,.12);border-radius:13px;padding:10px;background:rgba(255,255,255,.05)}.richCard b{display:block;color:#ffe4a3}.richCard span{display:block;color:#dce7f6;font-size:12px;margin-top:5px}
    .richTable{display:grid;gap:6px}.richRow{display:grid;grid-template-columns:160px 1fr;gap:8px}.richRow b{color:#ffe4a3}.richRow span{color:#dce7f6}
    .docReport{display:block;font-family:Inter,ui-sans-serif,system-ui,sans-serif;color:#eef7ff}
    .docHero{display:block;position:relative;overflow:hidden;border:1px solid rgba(97,211,148,.38);border-radius:18px;padding:16px;margin-bottom:12px;background:radial-gradient(circle at 90% 0%,rgba(97,211,148,.24),transparent 34%),linear-gradient(135deg,rgba(97,211,148,.15),rgba(112,215,255,.08))}
    .docHero::after{content:"";position:absolute;inset:-80% -40%;background:linear-gradient(90deg,transparent,rgba(255,255,255,.12),transparent);animation:sweep 4s linear infinite}
    .docHero>*{position:relative;z-index:1}
    .docFlow{display:grid;grid-template-columns:repeat(5,minmax(0,1fr));align-items:center;gap:8px;margin:10px 0 12px}
    .docNode{display:grid;place-items:center;min-height:70px;border:1px solid rgba(112,215,255,.26);border-radius:15px;background:linear-gradient(180deg,rgba(9,24,39,.95),rgba(2,9,18,.95));text-align:center;box-shadow:inset 0 1px 0 rgba(255,255,255,.08)}
    .docNode b{display:block;color:#fff;font-size:14px}.docNode span{display:block;color:#b8c5d8;font-size:11px;margin-top:3px}.docArrow{height:3px;background:linear-gradient(90deg,#70d7ff,#61d394);border-radius:999px;position:relative;overflow:hidden}.docArrow::after{content:"";position:absolute;inset:0;background:linear-gradient(90deg,transparent,#fff,transparent);animation:sweep 1.2s linear infinite}
    .docCards{display:grid;grid-template-columns:repeat(4,minmax(0,1fr));gap:8px;margin:0 0 12px}.docCard{display:block;border:1px solid rgba(255,255,255,.12);border-radius:14px;padding:10px;background:rgba(255,255,255,.055)}.docCard b{display:block;color:#f8d998;font-size:12px;text-transform:uppercase;letter-spacing:.06em}.docCard strong{display:block;color:#fff;font-size:20px;line-height:1.08;margin-top:4px}.docCard span{display:block;color:#cbd8e9;font-size:11px;margin-top:5px}
    .docSection{display:block;border:1px solid rgba(112,215,255,.24);border-radius:14px;padding:12px;margin-bottom:10px;background:rgba(255,255,255,.04)}.docSection h3{margin:0 0 8px;color:#70e39d;font-size:16px}.docSection p{color:#dce7f6;font-size:13px;line-height:1.38}.docSection ul{margin:7px 0 0 18px;padding:0}.docSection li{margin:5px 0;color:#dce7f6;font-size:13px;line-height:1.3}.docMeta{display:grid;grid-template-columns:150px 1fr;gap:6px}.docMeta b{color:#f8d998}.docMeta span{color:#dce7f6}
    .heroDiagram{height:100%;display:grid;place-items:center;border:1px solid rgba(255,255,255,.1);border-radius:14px;background:radial-gradient(circle at 50% 40%,rgba(97,211,148,.14),transparent 34%),rgba(0,0,0,.16);overflow:hidden}
    svg{width:100%;height:100%}.node{fill:#121b2b;stroke:#4b95bd;stroke-width:3}.node.hot{stroke:#61d394;filter:drop-shadow(0 0 16px rgba(97,211,148,.32))}.wire{fill:none;stroke:#70d7ff;stroke-width:4;stroke-dasharray:12 12;animation:dash 1.3s linear infinite}.wire.green{stroke:#61d394}.nodeText{fill:#f7fbff;font-size:20px;font-weight:900}.subText{fill:#bcc9dc;font-size:14px;font-weight:800}
    .watermark{position:fixed;right:18px;bottom:10px;color:rgba(247,251,255,.38);font-size:10px;font-weight:850;letter-spacing:.08em;text-transform:uppercase;pointer-events:none}
    @keyframes pulse{50%{transform:translateY(-1px);box-shadow:0 0 0 3px rgba(97,211,148,.16),0 0 42px rgba(97,211,148,.3)}}@keyframes sweep{to{transform:translateX(70%)}}@keyframes clickpop{0%{transform:scale(1)}28%{transform:scale(.96);box-shadow:0 0 0 3px rgba(248,217,152,.25),0 0 36px rgba(248,217,152,.22)}100%{transform:scale(1)}}@keyframes ripple{0%{opacity:.9;transform:scale(.25)}100%{opacity:0;transform:scale(1.8)}}@keyframes blink{50%{opacity:0}}@keyframes dash{to{stroke-dashoffset:-48}}
    @media(max-width:1120px){header{grid-template-columns:1fr;align-items:stretch}.status{justify-self:start}.inputs{grid-template-columns:1fr 1fr 1fr}}
    @media(max-height:760px){main{padding:10px;gap:8px}.captionBar{min-height:74px}.captionText{font-size:21px}.explain p{font-size:15px}.logo{width:42px;height:42px}h1{font-size:34px}.richTitle{font-size:21px}pre{font-size:13px;line-height:1.5}.flow{display:none}}
  </style>
</head>
<body>
<publishing-demo>
  <main>
    <header>
      <div class="brand">
        <a href="/demo" title="Open landing page">
          <img class="logo" src="/img/DeltaSignalIcon-256.png" alt="DeltaSignal logo">
          <div><div class="eyebrow">Google Cloud Rapid Agent Hackathon · MongoDB Track</div><h1>AITrailblazer AI Agent Publishing</h1></div>
        </a>
      </div>
      <div class="inputs"><div><label>Private demo key</label><input id="key" type="password" autocomplete="off" placeholder="PUBLISHING_DEMO_API_KEY"></div><div><label>TripCode</label><input id="tripcode" value="HUT-RIVER-001"></div><div><label>Session</label><input id="session" value="publisher-demo"></div></div>
      <div class="status"><span class="dot"></span> Cloud Run · Gemini · MongoDB MCP · TripCodes</div>
    </header>
    <section class="grid">
      <aside class="panel control">
        <div class="panelHead"><div><h2>3 Minute Publishing Proof</h2><p>One article archive becomes agent memory.</p></div><span class="pill">Live UI</span></div>
        <div class="metrics"><div class="metric"><strong>300+</strong><span>source articles</span></div><div class="metric"><strong>3</strong><span>River nodes</span></div><div class="metric"><strong>MCP</strong><span>database access</span></div></div>
        <div class="flow"><div><b>Brief</b>Archive</div><div><b>Resolve</b>TripCode</div><div><b>Proof</b>Runtime</div><div><b>Memory</b>Follow-up</div><div><b>Usage</b>Ledger</div></div>
        <div class="buttons">
          <button data-action="overview"><span class="btnIcon"><svg viewBox="0 0 24 24"><path d="M4 5h16M4 12h10M4 19h16"/></svg></span><span class="btnLabel">Overview</span></button><button data-action="diagram"><span class="btnIcon"><svg viewBox="0 0 24 24"><path d="M6 7h5v5H6zM13 12h5v5h-5zM11 9h4M9 12v3"/></svg></span><span class="btnLabel">Diagram</span></button><button data-action="archive"><span class="btnIcon"><svg viewBox="0 0 24 24"><path d="M5 6h14v14H5zM8 6V4h8v2M8 11h8M8 15h5"/></svg></span><span class="btnLabel">Archive</span></button><button data-action="tripcode"><span class="btnIcon"><svg viewBox="0 0 24 24"><path d="M7 7h4v4H7zM13 7h4v4h-4zM7 13h4v4H7zM13 13h1M17 13h1M13 17h5"/></svg></span><span class="btnLabel">TripCode</span></button>
          <button data-action="judge"><span class="btnIcon"><svg viewBox="0 0 24 24"><path d="M12 3l7 4v5c0 4-3 7-7 9-4-2-7-5-7-9V7zM9 12l2 2 4-5"/></svg></span><span class="btnLabel">Judge Proof</span></button><button data-action="render"><span class="btnIcon"><svg viewBox="0 0 24 24"><path d="M4 5h16v14H4zM7 9h10M7 13h6M15 13h2"/></svg></span><span class="btnLabel">Render</span></button><button data-action="followup"><span class="btnIcon"><svg viewBox="0 0 24 24"><path d="M7 8a5 5 0 0 1 8-3l2 2M17 16a5 5 0 0 1-8 3l-2-2M8 17H4v4M16 7h4V3"/></svg></span><span class="btnLabel">Memory</span></button><button data-action="usage"><span class="btnIcon"><svg viewBox="0 0 24 24"><path d="M4 19V5M4 19h16M8 16v-5M12 16V8M16 16v-3"/></svg></span><span class="btnLabel">Usage</span></button>
        </div>
        <div class="explain"><h3 id="actionTitle">Ready</h3><p id="actionText">Click Overview, or open this page with <code>?autoplay=1</code> to run the full three-minute sequence.</p></div>
      </aside>
      <section class="panel workspace">
        <div class="outHead"><div><h2>Response Workspace</h2><p>Request trace, JSON proof, rendered packet, and session memory.</p></div><div class="tools"><span class="pill" id="host">local</span><button data-scroll="top" class="secondary">Top</button><button data-scroll="page" class="secondary">Scroll</button><button data-scroll="end" class="secondary">End</button></div></div>
        <pre id="out">Ready. The autoplay sequence will show archive brief, TripCode/River resolution, judge proof, memory follow-up, and usage.</pre>
      </section>
    </section>
    <section class="captionBar" aria-label="Demo explanation"><span id="captionText" class="captionText"></span></section>
  </main>
  <div class="watermark">© 2026 AITrailblazer</div>
</publishing-demo>
<script>
class PublishingDemo extends HTMLElement{
  connectedCallback(){queueMicrotask(()=>this.bind())}
  bind(){if(this.bound)return;this.out=this.querySelector('#out');this.key=this.querySelector('#key');this.tripcode=this.querySelector('#tripcode');this.session=this.querySelector('#session');this.actionTitle=this.querySelector('#actionTitle');this.actionText=this.querySelector('#actionText');this.captionText=this.querySelector('#captionText');if(!this.out||!this.key||!this.tripcode||!this.session||!this.actionTitle||!this.actionText||!this.captionText)return;this.bound=true;this.querySelector('#host').textContent=location.host||'Cloud Run';this.autoplay=new URLSearchParams(location.search).get('autoplay')==='1';for(const b of this.querySelectorAll('button[data-action]'))b.addEventListener('click',()=>this.run(b.dataset.action));for(const b of this.querySelectorAll('button[data-scroll]'))b.addEventListener('click',()=>this.scroll(b.dataset.scroll,b));this.describe('overview');this.showOverview();if(this.autoplay)setTimeout(()=>this.play(),1200)}
  headers(json=false){const h={};if(json)h['Content-Type']='application/json';const key=this.key.value.trim();if(key)h['X-Demo-Key']=key;return h}
  copy(action){return {overview:{title:'Overview',text:'This demo proves a generalized publishing agent: article archives become TripCodes, River memory, MongoDB-backed context, Gemini synthesis, and follow-up session state.',caption:'On screen: the control panel introduces the archive-to-agent-memory workflow, while the response workspace lists the full three-minute path.',scrollCaption:'The overview explains the full flow before protected requests run.'},diagram:{title:'Architecture Diagram',text:'The user starts with archive intent. MongoDB memory, TripCodes, River edges, and MCP context feed the Cloud Run publishing agent, which returns a grounded reader packet.',caption:'On screen: the diagram shows inputs flowing into the publishing agent and a grounded output packet flowing back to the user.',scrollCaption:'The visual proof shows that the agent is not a chatbot over text; it is an orchestrated memory and retrieval workflow.'},archive:{title:'Archive Brief',text:'The browser posts a publication archive request. The agent returns a plan for turning a DeltaSignal-style archive into article objects, TripCodes, Rivers, vector memory, and reader prompts.',caption:'On screen: the live JSON response shows the archive plan, findings, disclosures, and cost metadata.',scrollCaption:'Now scrolling the archive brief: watch how the agent turns publication inventory into TripCodes, Rivers, article objects, vector memory, and reader prompts.'},tripcode:{title:'TripCode Resolve',text:'The protected route resolves HUT-RIVER-001 into a publication object, anchor article, three-node River, MongoDB memory fit, execution trace, and session memory.',caption:'On screen: TripCode resolution returns publication identity, article claims, River nodes, MongoDB collections, monitor-next prompts, and disclosures.',scrollCaption:'Now scrolling the TripCode result: article identity, claims, River nodes, MongoDB fit, execution trace, session memory, and boundaries stay visible.'},judge:{title:'Judge Proof Package',text:'The judge proof endpoint packages the runtime proof: required panels, required collections, live flows, TripCode result, follow-up memory, and runtime integration status.',caption:'On screen: the proof package shows the entire competition workflow in one response.',scrollCaption:'Now scrolling judge proof: the packaged evidence shows panels, collections, live flows, runtime proof, TripCode result, and memory follow-up in one artifact.'},render:{title:'Rendered Packet',text:'The latest JSON proof is converted into a rich readable page. This is what a first-time user can inspect quickly after the raw request returns.',caption:'On screen: the raw proof is transformed into a human-readable evidence packet with summary, identity, sources, memory, and boundaries.',scrollCaption:'Now scrolling the rendered packet: summary, workflow identity, River memory, runtime proof, usage, and raw evidence are readable without parsing JSON.'},followup:{title:'Second-turn Memory',text:'The user asks a follow-up without repeating the TripCode. The agent uses stored session memory so the River remains in context.',caption:'On screen: the follow-up response proves session memory by returning prior TripCode and River context from the same session.',scrollCaption:'Now scrolling memory: the response shows session-memory mode, the prior TripCode, remembered River context, and what should stay in context.'},usage:{title:'Usage Ledger',text:'The usage endpoint shows cost-aware request accounting and remaining budget metadata so the demo stays measurable.',caption:'On screen: the usage response shows local-estimate cost controls for the demo run.',scrollCaption:'Now scrolling usage: request kinds, local estimates, remaining budget, and source notes remain visible for cost-aware operation.'}}[action]||{title:'Action',text:'The selected workflow runs and returns bounded agent output.',caption:'The selected workflow returned.'}}
  async describe(action){const c=this.copy(action);this.current=c;this.actionTitle.textContent=c.title;this.captionText.innerHTML='';await this.type(this.actionText,c.text,this.autoplay?58:68,'actionCursor');await this.type(this.captionText,c.caption||c.text,this.autoplay?82:92,'cursor')}
  async type(el,text,delay,cls){el.innerHTML='';let s='';for(const word of text.split(' ')){s+=(s?' ':'')+word;el.innerHTML=this.esc(s)+'<span class="'+cls+'"></span>';await new Promise(r=>setTimeout(r,delay))}el.textContent=text}
  async request(label,fn,button){this.setRunning(button);const started=new Date().toISOString();let timer;let tick=0;const update=()=>{tick++;this.out.innerHTML=this.trace(label,'running',started,(tick/10).toFixed(1)+'s')};update();timer=setInterval(update,100);try{const t0=performance.now();const res=await fn();const text=await res.text();const ms=Math.round(performance.now()-t0)+'ms';let body=text,parsed=null,highlight='';try{parsed=JSON.parse(text);body=JSON.stringify(parsed,null,2);highlight=this.color(body)}catch{highlight=this.esc(body)}clearInterval(timer);this.last={label,status:res.status,ok:res.ok,started,duration:ms,body,parsed,highlight,header:'<span class="resp-title">'+this.esc(label)+'</span><br><span class="'+(res.ok?'resp-status':'resp-error')+'">HTTP '+res.status+' · measured network time '+ms+'</span><br><br>',trace:this.trace(label,'completed',started,ms,res.status)};await this.stream(this.last.trace,this.last.header,highlight)}catch(e){clearInterval(timer);this.out.innerHTML='<span class="resp-title">'+this.esc(label)+'</span><br><span class="resp-error">ERROR</span><br><br>'+this.esc(e&&e.message?e.message:String(e))}finally{this.setRunning(null)}}
  showOverview(){this.out.innerHTML='<span class="rich"><span class="richHero"><span class="richKicker">Three-minute proof</span><span class="richTitle">Archive memory becomes an agent workflow.</span><span class="richSummary">A reader starts with a TripCode from a publication. The agent resolves article records, River edges, MongoDB memory fit, Gemini-ready synthesis, session continuity, and usage metadata.</span></span><span class="richGrid"><span class="richCard"><b>Archive</b><span>300+ article publication proof case.</span></span><span class="richCard"><b>TripCode</b><span>HUT-RIVER-001 resolves the anchor article.</span></span><span class="richCard"><b>Memory</b><span>Second-turn follow-up reuses River context.</span></span></span><span class="richSection"><h3>Autoplay path</h3><span class="richTable"><span class="richRow"><b>1</b><span>Archive brief</span></span><span class="richRow"><b>2</b><span>TripCode/River resolve</span></span><span class="richRow"><b>3</b><span>Judge proof package</span></span><span class="richRow"><b>4</b><span>Rendered proof packet</span></span><span class="richRow"><b>5</b><span>Session memory follow-up</span></span><span class="richRow"><b>6</b><span>Usage ledger</span></span></span></span></span>';this.out.scrollTop=0}
  showDiagram(){this.out.innerHTML='<span class="heroDiagram"><svg viewBox="0 0 920 520" role="img" aria-label="Archive memory architecture"><defs><marker id="a" markerWidth="14" markerHeight="14" refX="12" refY="7" orient="auto"><path d="M0,0 L14,7 L0,14 Z" fill="#70d7ff"/></marker><marker id="g" markerWidth="14" markerHeight="14" refX="12" refY="7" orient="auto"><path d="M0,0 L14,7 L0,14 Z" fill="#61d394"/></marker></defs><path class="wire" d="M162 132 H342" marker-end="url(#a)"/><path class="wire green" d="M730 132 H596" marker-end="url(#g)"/><path class="wire green" d="M280 362 C350 292 440 276 460 210" marker-end="url(#g)"/><path class="wire green" d="M640 362 C570 292 480 276 460 210" marker-end="url(#g)"/><path class="wire" d="M460 210 V266" marker-end="url(#a)"/><rect class="node hot" x="48" y="72" width="150" height="106" rx="24"/><text class="nodeText" x="98" y="120">User</text><text class="subText" x="82" y="147">TripCode ask</text><rect class="node hot" x="342" y="72" width="236" height="138" rx="26"/><text class="nodeText" x="408" y="128">Cloud Run</text><text class="subText" x="396" y="158">Publishing agent</text><rect class="node" x="712" y="72" width="160" height="106" rx="24"/><text class="nodeText" x="750" y="120">MCP</text><text class="subText" x="736" y="147">MongoDB context</text><rect class="node hot" x="360" y="270" width="200" height="78" rx="22"/><text class="nodeText" x="410" y="306">Output</text><text class="subText" x="392" y="330">reader packet</text><rect class="node" x="160" y="360" width="240" height="96" rx="24"/><text class="nodeText" x="226" y="404">River</text><text class="subText" x="200" y="430">article continuity</text><rect class="node" x="560" y="360" width="180" height="96" rx="24"/><text class="nodeText" x="594" y="404">Archive</text><text class="subText" x="588" y="430">article memory</text></svg></span>'}
  async play(){const seq=[['diagram',5500],['archive',4500],['tripcode',3500],['judge',3500],['render',4500],['followup',3500],['usage',4500]];for(const [a,p]of seq){const b=this.querySelector('button[data-action="'+a+'"]');if(b){b.classList.add('running');setTimeout(()=>b.classList.remove('running'),900)}await this.run(a);await new Promise(r=>setTimeout(r,p))}this.type(this.captionText,'Demo complete: the publishing agent turned one archive and TripCode into a grounded packet, River memory, follow-up context, and usage visibility.',82,'cursor')}
  async run(action){const b=this.querySelector('button[data-action="'+action+'"]');this.flash(b);const d=this.describe(action);const q='How can this archive become agent-readable memory for other publishers?';const map={overview:()=>{this.setRunning(b);this.showOverview();setTimeout(()=>this.setRunning(null),620)},diagram:()=>{this.setRunning(b);this.showDiagram();setTimeout(()=>this.setRunning(null),620)},archive:()=>this.request('POST /v1/archive-brief',()=>fetch('/v1/archive-brief',{method:'POST',headers:this.headers(true),body:JSON.stringify({publication_url:'https://deltasignal.substack.com/',publication:'DeltaSignal',question:q})}),b),tripcode:()=>this.request('POST /v1/tripcode',()=>fetch('/v1/tripcode',{method:'POST',headers:this.headers(true),body:JSON.stringify({tripcode:this.tripcode.value.trim(),session_id:this.session.value.trim(),include_river:true,include_mongo_fit:true,question:'Resolve this TripCode into article memory, River context, MongoDB fit, and monitor-next actions.'})}),b),judge:()=>this.request('GET /v1/judge-demo',()=>fetch('/v1/judge-demo',{headers:this.headers()}),b),render:()=>{this.setRunning(b);this.render();setTimeout(()=>this.setRunning(null),620)},followup:()=>this.request('POST /resolve session follow-up',()=>fetch('/resolve?session_id='+encodeURIComponent(this.session.value.trim()),{method:'POST',headers:this.headers(),body:'Using the previous publication River, what should stay in context?'}),b),usage:()=>this.request('GET /v1/usage',()=>fetch('/v1/usage',{headers:this.headers()}),b)};const result=await map[action]();await d;await this.inspect();return result}
  render(){if(!this.last||!this.last.parsed){this.out.innerHTML='<span class="resp-title">Rendered Packet</span><br><span class="resp-error">Run a JSON step first.</span>';return}const data=this.last.parsed;const trip=this.read(data.tripcode)||this.read(data.tripcode_result&&data.tripcode_result.tripcode)||this.tripcode.value;const mode=this.read(data.mode)||'response';const summary=this.read(data.summary)||this.read(data.brief)||this.read(data.one_line)||this.read(data.answer)||'Structured publishing-agent response.';const nodes=this.read(data.packet&&data.packet.river&&data.packet.river.node_count)||this.read(data.tripcode_result&&data.tripcode_result.packet&&data.tripcode_result.packet.river&&data.tripcode_result.packet.river.node_count)||this.read(data.dataset_boundary&&data.dataset_boundary.articles)||'available';const memory=this.read(data.memory&&data.memory.available)||this.read(data.follow_up_result&&data.follow_up_result.memory&&data.follow_up_result.memory.available)||'tracked';const cost=this.read(data.cost&&data.cost.estimated_total_usd)||this.read(data.cost&&data.cost.estimated_request_usd)||'visible';this.out.innerHTML='<span class="docReport"><span class="docHero"><span class="richKicker">Rendered HTML report from live JSON</span><span class="richTitle">'+this.esc(this.last.label)+'</span><span class="richSummary">'+this.esc(summary)+'</span></span><span class="docFlow"><span class="docNode"><b>Archive</b><span>publication objects</span></span><span class="docArrow"></span><span class="docNode"><b>TripCode</b><span>'+this.esc(trip)+'</span></span><span class="docArrow"></span><span class="docNode"><b>Agent Packet</b><span>memory + next actions</span></span></span><span class="docCards"><span class="docCard"><b>Mode</b><strong>'+this.esc(mode)+'</strong><span>runtime response type</span></span><span class="docCard"><b>TripCode</b><strong>'+this.esc(trip)+'</strong><span>stable resolver handle</span></span><span class="docCard"><b>River</b><strong>'+this.esc(nodes)+'</strong><span>continuity signal</span></span><span class="docCard"><b>Cost</b><strong>'+this.esc(cost)+'</strong><span>usage visibility</span></span></span><span class="docSection"><h3>What The Agent Did</h3><p>'+this.esc(summary)+'</p><ul>'+this.listItems(data.required_panels||data.live_flows||data.findings||data.actions||[]).join('')+'</ul></span>'+this.docSection('Execution Trace',data.execution_trace||data.runtime_proof||{route:this.last.label,http_status:this.last.status,duration:this.last.duration})+this.docSection('Archive / River Memory',data.packet||data.tripcode_result||data.required_collections)+this.docSection('Session Memory',data.memory||data.follow_up_result||{available:memory})+this.docSection('Evidence Boundaries',data.disclosures||data.submission_cuts||data.dataset_boundary)+this.docSection('Usage',data.cost||data.usage||{cost:cost})+'</span>';this.out.scrollTop=0}
  docSection(t,v){const rows=Object.entries(this.flatten(v)).slice(0,10).map(([k,x])=>'<span class="docMeta"><b>'+this.esc(k)+'</b><span>'+this.esc(this.read(x))+'</span></span>').join('');return '<span class="docSection"><h3>'+this.esc(t)+'</h3>'+rows+'</span>'}
  listItems(v){if(!Array.isArray(v))return [];return v.slice(0,5).map(x=>'<li>'+this.esc(this.read(x))+'</li>')}
  section(t,v){return '<span class="richSection"><h3>'+this.esc(t)+'</h3><span class="richTable">'+Object.entries(this.flatten(v)).slice(0,12).map(([k,x])=>'<span class="richRow"><b>'+this.esc(k)+'</b><span>'+this.esc(this.read(x))+'</span></span>').join('')+'</span></span>'}
  flatten(v){if(!v)return {status:'not returned'};if(Array.isArray(v))return Object.fromEntries(v.slice(0,10).map((x,i)=>[String(i+1),x]));if(typeof v==='object')return v;return {value:v}}
  async inspect(){const max=this.out.scrollHeight-this.out.clientHeight;if(max<80||!this.current.scrollCaption)return;this.out.scrollTop=0;await new Promise(r=>setTimeout(r,520));const c=this.type(this.captionText,this.current.scrollCaption,this.autoplay?82:92,'cursor');await new Promise(r=>setTimeout(r,520));await this.animate(0,max,Math.min(this.autoplay?24000:12800,Math.max(this.autoplay?11000:6200,max*(this.autoplay?12.0:8.2))));await c}
  async stream(trace,header,html){this.out.innerHTML=trace+header+'<span id="st"></span><span class="cursor"></span>';const target=this.querySelector('#st');let s='';for(const ch of this.chunks(html,105)){s+=ch;target.innerHTML=s;this.out.scrollTop=0;await new Promise(r=>setTimeout(r,16))}this.out.innerHTML=trace+header+html;this.out.scrollTop=0}
  trace(label,state,started,elapsed,status){return '<span class="trace"><b>route</b> <span class="route">'+this.esc(label)+'</span><br><b>state</b> <span class="'+(state==='completed'?'ok':'')+'">'+this.esc(state)+'</span><br><b>started</b> '+this.esc(started)+'<br><b>elapsed</b> '+this.esc(elapsed)+(status?'<br><b>status</b> <span class="ok">HTTP '+this.esc(status)+'</span>':'')+'</span><br>'}
  scroll(mode,b){b.classList.add('running');const start=this.out.scrollTop;let target=0;if(mode==='page')target=Math.min(this.out.scrollHeight-this.out.clientHeight,start+Math.round(this.out.clientHeight*.78));if(mode==='end')target=this.out.scrollHeight-this.out.clientHeight;this.animate(start,Math.max(0,target),900).finally(()=>b.classList.remove('running'))}
  async animate(start,target,duration){const t0=Date.now();while(true){const t=Math.min(1,(Date.now()-t0)/duration);const e=t<.5?2*t*t:1-Math.pow(-2*t+2,2)/2;this.out.scrollTop=start+(target-start)*e;if(t>=1)break;await new Promise(r=>setTimeout(r,16))}}
  setRunning(b){for(const x of this.querySelectorAll('button.running'))x.classList.remove('running');if(b)b.classList.add('running')}
  flash(b){if(!b)return;b.classList.remove('clicked');void b.offsetWidth;b.classList.add('clicked');setTimeout(()=>b.classList.remove('clicked'),680)}
  color(json){let out='',i=0;while(i<json.length){const ch=json[i];if(ch==='"'){const st=i;i++;let esc=false;while(i<json.length){const c=json[i];if(esc)esc=false;else if(c==='\\')esc=true;else if(c==='"'){i++;break}i++}let j=i;while(j<json.length&&/\s/.test(json[j]))j++;out+='<span class="'+(json[j]===':'?'json-key':'json-string')+'">'+this.esc(json.slice(st,i))+'</span>';continue}const m=json.slice(i).match(/^-?\d+(?:\.\d+)?/);if(m){out+='<span class="json-number">'+m[0]+'</span>';i+=m[0].length;continue}if(json.startsWith('true',i)||json.startsWith('false',i)){const v=json.startsWith('true',i)?'true':'false';out+='<span class="json-bool">'+v+'</span>';i+=v.length;continue}if(json.startsWith('null',i)){out+='<span class="json-null">null</span>';i+=4;continue}if('{}[]'.includes(ch)){out+='<span class="json-brace">'+ch+'</span>';i++;continue}if(ch===':'){out+='<span class="json-colon">:</span>';i++;continue}if(ch===','){out+='<span class="json-punc">,</span>';i++;continue}out+=this.esc(ch);i++}return out}
  chunks(html,n){const out=[];let cur='',vis=0;for(let i=0;i<html.length;){if(html[i]==='<'){const e=html.indexOf('>',i);if(e<0)break;cur+=html.slice(i,e+1);i=e+1;continue}cur+=html[i++];vis++;if(vis>=n||html[i-1]==='\n'){out.push(cur);cur='';vis=0}}if(cur)out.push(cur);return out}
  read(v){if(v==null)return '';if(typeof v==='string'||typeof v==='number'||typeof v==='boolean')return String(v);if(Array.isArray(v))return v.map(x=>this.read(x)).filter(Boolean).join(' · ');if(typeof v==='object')return Object.entries(v).slice(0,8).map(([k,x])=>k+': '+this.read(x)).join(' · ');return String(v)}
  esc(v){return String(v).replace(/[&<>"']/g,ch=>({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[ch]))}
}
customElements.define('publishing-demo',PublishingDemo);
</script>
</body>
</html>`

func registerDemoUIRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /demo", servePublishingDemoLanding)
	mux.HandleFunc("GET /demo/", servePublishingDemoLanding)
	mux.HandleFunc("GET /demo/run", servePublishingDemoUI)
	mux.HandleFunc("GET /demo/run/", servePublishingDemoUI)
}

func servePublishingDemoLanding(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(publishingDemoLandingHTML))
}

func servePublishingDemoUI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	body := publishingDemoUIHTML
	if os.Getenv("PUBLISHING_DEMO_API_KEY") == "local-demo-key" {
		body = strings.Replace(body, `id="key" type="password" autocomplete="off" placeholder="PUBLISHING_DEMO_API_KEY"`, `id="key" type="password" autocomplete="off" placeholder="PUBLISHING_DEMO_API_KEY" value="local-demo-key"`, 1)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(body))
}
