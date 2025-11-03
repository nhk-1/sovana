import PropTypes from 'prop-types';

/**
 * Responsive container enforcing the design max-width.
 */
const Container = ({ as: Component = 'div', children, className = '', id }) => (
  <Component id={id} className={`mx-auto w-full max-w-content px-4 sm:px-6 lg:px-8 ${className}`}>
    {children}
  </Component>
);

Container.propTypes = {
  as: PropTypes.elementType,
  children: PropTypes.node.isRequired,
  className: PropTypes.string,
  id: PropTypes.string,
};

export default Container;
