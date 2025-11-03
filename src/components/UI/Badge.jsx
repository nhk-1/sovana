import PropTypes from 'prop-types';

/**
 * Pill-shaped badge for highlighting short text.
 */
const Badge = ({ children }) => (
  <span className="inline-flex items-center rounded-full border border-ocean/30 bg-ocean/10 px-3 py-1 text-xs font-semibold uppercase tracking-wide text-ocean dark:border-white/30 dark:bg-white/10 dark:text-white">
    {children}
  </span>
);

Badge.propTypes = {
  children: PropTypes.node.isRequired,
};

export default Badge;
