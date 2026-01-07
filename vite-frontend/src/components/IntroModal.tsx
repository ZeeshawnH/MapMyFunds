type IntroModalProps = {
  onClose: () => void;
};

const IntroModal = ({ onClose }: IntroModalProps) => {
  return (
    <div className="intro-modal-backdrop" role="dialog" aria-modal="true">
      <div className="intro-modal">
        <div className="intro-modal-brand">Map My Funds</div>
        <h2 className="intro-modal-title">
          How Each State Voted With Their Money
        </h2>
        <p className="intro-modal-text">
          This project visualizes federal presidential campaign contributions by
          state and candidate. It uses public data from the Federal Election
          Commission, aggregated by election cycle. The map is styled like an
          electoral map: each state is color-coded by the party of the top
          fundraising candidate in that state.
        </p>
        <p className="intro-modal-text">
          Use the year selector above the map to switch between election cycles.
          Hover over a state to see which campaigns raised the most money there,
          and explore the sidebars for top contributing states and candidates.
        </p>
        <p className="intro-modal-note">
          Learn more{" "}
          <a
            href="https://github.com/ZeeshawnH/MapMyFunds"
            target="_blank"
            rel="noopener noreferrer"
          >
            here
          </a>
          .
        </p>
        <button type="button" className="intro-modal-button" onClick={onClose}>
          Got it
        </button>
      </div>
    </div>
  );
};

export default IntroModal;
