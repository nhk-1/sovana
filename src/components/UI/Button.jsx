import PropTypes from 'prop-types';
import { motion } from 'framer-motion';
import { Link as ScrollLink } from 'react-scroll';

/**
 * Accessible, reusable button component supporting motion and anchor scroll.
 */
const Button = ({ children, href, onClick, variant = 'primary', className = '', as = 'button' }) => {
  const baseClasses =
    'inline-flex items-center justify-center rounded-full px-6 py-3 text-sm font-semibold transition focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2';
  const variants = {
    primary: 'bg-ocean text-white shadow-subtle hover:bg-midnight focus-visible:outline-ocean',
    secondary:
      'border border-ocean/40 bg-transparent text-ocean hover:border-ocean focus-visible:outline-ocean dark:border-white/50 dark:text-white dark:hover:border-white',
  };

  if (href) {
    if (href.startsWith('#')) {
      return (
        <ScrollLink
          to={href.replace('#', '')}
          smooth
          duration={600}
          offset={-80}
          className={`${baseClasses} cursor-pointer ${variants[variant]} ${className}`}
        >
          {children}
        </ScrollLink>
      );
    }

    return (
      <a href={href} className={`${baseClasses} ${variants[variant]} ${className}`} onClick={onClick}>
        {children}
      </a>
    );
  }

  const Component = motion[as] || motion.button;

  return (
    <Component
      type={as === 'button' ? 'button' : undefined}
      whileHover={{ y: -2 }}
      whileTap={{ scale: 0.98 }}
      onClick={onClick}
      className={`${baseClasses} ${variants[variant]} ${className}`}
    >
      {children}
    </Component>
  );
};

Button.propTypes = {
  children: PropTypes.node.isRequired,
  href: PropTypes.string,
  onClick: PropTypes.func,
  variant: PropTypes.oneOf(['primary', 'secondary']),
  className: PropTypes.string,
  as: PropTypes.oneOf(['button', 'div']),
};

export default Button;
